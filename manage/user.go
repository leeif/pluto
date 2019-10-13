package manage

import (
	"encoding/json"
	"errors"
	"time"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/models"

	saltUtil "github.com/leeif/pluto/utils/salt"

	b64 "encoding/base64"

	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/mail"
)

func (m *Manger) ResetPasswordMail(rpm request.ResetPasswordMail, domain string) *perror.PlutoError {

	user := models.User{}
	identifyToken := b64.StdEncoding.EncodeToString([]byte(rpm.Mail))
	if m.db.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken).First(&user).RecordNotFound() {
		return perror.MailIsNotExsit
	}

	ml := mail.NewMail(m.config)
	if err := ml.SendResetPassword(*user.Mail, domain); err != nil {
		return err
	}

	return nil
}

func (m *Manger) ResetPasswordPage(token string) *perror.PlutoError {

	header, payload, err := jwt.VerifyB64JWT(token)
	// token verify failed
	if err != nil {
		return err
	}

	head := jwt.Head{}
	json.Unmarshal(header, &head)
	if head.Type != jwt.PASSWORDRESET {
		return perror.InvalidJWTToekn
	}

	prp := jwt.PasswordResetPayload{}
	json.Unmarshal(payload, &prp)

	if time.Now().Unix() > prp.Expire {
		return perror.InvalidJWTToekn
	}

	user := models.User{}
	identifyToken := b64.StdEncoding.EncodeToString([]byte(prp.Mail))
	if m.db.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken).First(&user).RecordNotFound() {
		return perror.ServerError.Wrapper(errors.New("mail not found"))
	}

	// user is updated after password reset token is created
	if user.UpdatedAt.Unix() > prp.Create {
		return perror.InvalidJWTToekn
	}

	return nil
}

func (m *Manger) ResetPassword(rp request.ResetPassword) *perror.PlutoError {

	header, payload, perr := jwt.VerifyB64JWT(rp.Token)
	if perr != nil {
		return perr
	}

	head := jwt.Head{}
	json.Unmarshal(header, &head)

	if head.Type != jwt.PASSWORDRESET {
		return perror.InvalidJWTToekn
	}

	prp := jwt.PasswordResetPayload{}
	json.Unmarshal(payload, &prp)

	if time.Now().Unix() > prp.Expire {
		return perror.InvalidJWTToekn
	}

	tx := m.db.Begin()
	defer func() {
		tx.Rollback()
	}()

	user := models.User{}
	identifyToken := b64.StdEncoding.EncodeToString([]byte(prp.Mail))
	if tx.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken).First(&user).RecordNotFound() {
		return perror.ServerError.Wrapper(errors.New("mail not found"))
	}

	// user is updated after password reset token is created
	if user.UpdatedAt.Unix() > prp.Create {
		return perror.InvalidJWTToekn
	}

	salt := models.Salt{}
	if tx.Where("user_id = ?", user.ID).First(&salt).RecordNotFound() {
		return perror.ServerError.Wrapper(errors.New("Salt is not found"))
	}

	salt.Salt = saltUtil.RandomSalt(prp.Mail)

	if err := update(tx, &salt); err != nil {
		return err
	}

	encodedPassword, err := saltUtil.EncodePassword(rp.Password, salt.Salt)
	if err != nil {
		return perror.ServerError.Wrapper(errors.New("Salt encoding is failed"))
	}

	user.Password = &encodedPassword
	if err := update(tx, &user); err != nil {
		return err
	}

	// add operation history
	if err := historyOperation(tx, OperationResetPassword, user.ID); err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (m *Manger) UserInfo(token string) (*models.User, *perror.PlutoError) {
	header, payload, err := jwt.VerifyB64JWT(token)
	if err != nil {
		return nil, err
	}

	head := jwt.Head{}
	json.Unmarshal(header, &head)

	if head.Type != jwt.ACCESS {
		return nil, perror.InvalidJWTToekn
	}

	userPayload := jwt.UserPayload{}
	json.Unmarshal(payload, &userPayload)

	if time.Now().Unix() > userPayload.Expire {
		return nil, perror.InvalidJWTToekn
	}

	user := models.User{}
	if m.db.Where("id = ?", userPayload.UserID).First(&user).RecordNotFound() {
		return nil, perror.ServerError.Wrapper(errors.New("user not found id: " + string(userPayload.UserID)))
	}

	return &user, nil
}

func (m *Manger) RefreshAccessToken(rat request.RefreshAccessToken) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)

	tx := m.db.Begin()

	defer func() {
		tx.Rollback()
	}()

	rt := models.RefreshToken{}
	if tx.Where("user_id = ? and refresh_token = ?", rat.UseID, rat.RefreshToken).First(&rt).RecordNotFound() {
		return nil, perror.InvalidRefreshToken
	}

	da := models.DeviceAPP{}
	da.ID = rt.DeviceAPPID
	if tx.Where("id = ?", da.ID).First(&da).RecordNotFound() {
		return nil, perror.InvalidRefreshToken
	}

	if rt.UserID != rat.UseID || da.DeviceID != rat.DeviceID || da.AppID != rat.AppID {
		return nil, perror.InvalidRefreshToken
	}

	// generate jwt token
	jwtToken, err := jwt.GenerateJWT(jwt.Head{Type: jwt.ACCESS},
		&jwt.UserPayload{UserID: rat.UseID, DeviceID: rat.DeviceID, AppID: rat.AppID}, m.config.JWT.AccessTokenExpire)

	if err != nil {
		return nil, err.Wrapper(errors.New("JWT token generate failed"))
	}

	// add operation history
	if err := historyOperation(tx, OperationRefreshToken, rat.UseID); err != nil {
		return nil, err
	}

	res["jwt"] = jwtToken

	tx.Commit()
	return res, nil
}
