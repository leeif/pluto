package manage

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/leeif/pluto/utils/avatar"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/models"

	saltUtil "github.com/leeif/pluto/utils/salt"

	b64 "encoding/base64"

	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/mail"
)

func (m *Manager) ResetPasswordMail(rpm request.ResetPasswordMail, baseURL string) *perror.PlutoError {

	user := models.User{}
	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(rpm.Mail))
	if m.db.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken).First(&user).RecordNotFound() {
		return perror.MailIsNotExsit
	}

	ml, err := mail.NewMail(m.config)
	if err != nil {
		return err
	}

	if err := ml.SendResetPassword(*user.Mail, baseURL); err != nil {
		return err
	}

	return nil
}

func (m *Manager) ResetPasswordPage(token string) *perror.PlutoError {

	jwtToken, err := jwt.VerifyB64JWT(token)
	// token verify failed
	if err != nil {
		return err
	}

	prp := jwt.PasswordResetPayload{}
	json.Unmarshal(jwtToken.Payload, &prp)

	if prp.Type != jwt.PASSWORDRESET {
		return perror.InvalidJWTToekn
	}

	if time.Now().Unix() > prp.Expire {
		return perror.JWTTokenExpired
	}

	user := models.User{}
	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(prp.Mail))
	if m.db.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken).First(&user).RecordNotFound() {
		return perror.ServerError.Wrapper(errors.New("mail not found"))
	}

	// user is updated after password reset token is created
	if user.UpdatedAt.Unix() > prp.Create {
		return perror.InvalidJWTToekn
	}

	return nil
}

func (m *Manager) ResetPassword(rp request.ResetPassword) *perror.PlutoError {

	jwtToken, perr := jwt.VerifyB64JWT(rp.Token)
	if perr != nil {
		return perr
	}

	prp := jwt.PasswordResetPayload{}
	json.Unmarshal(jwtToken.Payload, &prp)

	if prp.Type != jwt.PASSWORDRESET {
		return perror.InvalidJWTToekn
	}

	if time.Now().Unix() > prp.Expire {
		return perror.JWTTokenExpired
	}

	tx := m.db.Begin()
	defer func() {
		tx.Rollback()
	}()

	user := models.User{}
	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(prp.Mail))
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

func (m *Manager) UserInfo(token string) (*models.User, *perror.PlutoError) {
	jwtToken, err := jwt.VerifyB64JWT(token)
	if err != nil {
		return nil, err
	}

	userPayload := jwt.UserPayload{}
	json.Unmarshal(jwtToken.Payload, &userPayload)

	if userPayload.Type != jwt.ACCESS {
		return nil, perror.InvalidJWTToekn
	}

	if time.Now().Unix() > userPayload.Expire {
		return nil, perror.JWTTokenExpired
	}

	user := models.User{}
	if m.db.Where("id = ?", userPayload.UserID).First(&user).RecordNotFound() {
		return nil, perror.ServerError.Wrapper(errors.New("user not found id: " + string(userPayload.UserID)))
	}

	return &user, nil
}

func (m *Manager) RefreshAccessToken(rat request.RefreshAccessToken) (map[string]string, *perror.PlutoError) {
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

	user := models.User{}
	if tx.Where("id = ?", rat.UseID).First(&user).RecordNotFound() {
		return nil, perror.ServerError.Wrapper(fmt.Errorf("UserID not found: %d", rat.UseID))
	}

	// generate jwt token
	up := jwt.NewUserPayload(rat.UseID, rat.DeviceID, rat.AppID, user.LoginType, m.config.JWT.AccessTokenExpire)
	token, err := jwt.GenerateRSAJWT(up)

	if err != nil {
		return nil, err.Wrapper(errors.New("JWT token generate failed"))
	}

	// add operation history
	if err := historyOperation(tx, OperationRefreshToken, rat.UseID); err != nil {
		return nil, err
	}

	res["jwt"] = token.String()

	tx.Commit()
	return res, nil
}

func (m *Manager) UpdateUserInfo(token string, uui request.UpdateUserInfo) *perror.PlutoError {
	jwtToken, err := jwt.VerifyB64JWT(token)
	if err != nil {
		return err
	}

	userPayload := jwt.UserPayload{}
	json.Unmarshal(jwtToken.Payload, &userPayload)

	if userPayload.Type != jwt.ACCESS {
		return perror.InvalidJWTToekn
	}

	if time.Now().Unix() > userPayload.Expire {
		return perror.JWTTokenExpired
	}

	if userPayload.LoginType != MAILLOGIN {
		return perror.InvalidJWTToekn
	}

	tx := m.db.Begin()

	defer func() {
		tx.Rollback()
	}()

	user := models.User{}
	if tx.Where("id = ?", userPayload.UserID).First(&user).RecordNotFound() {
		return perror.ServerError.Wrapper(errors.New("user not found id: " + string(userPayload.UserID)))
	}

	if uui.Avatar != "" && m.isValidURL(uui.Avatar) {
		user.Avatar = uui.Avatar
	} else if uui.Avatar != "" && m.isValidBase64(uui.Avatar) {
		ag := avatar.AvatarGen{}
		ar, err := ag.GenFromBase64String(uui.Avatar)
		if err != nil {
			return err
		}
		as := avatar.NewAvatarSaver(m.config)
		url, err := as.SaveAvatarImageInOSS(ar)
		if err != nil {
			return err
		}
		user.Avatar = url
	} else if uui.Avatar != "" {
		return perror.InvalidAvatarFormat
	}

	if uui.Gender != "" {
		user.Gender = &uui.Gender
	}

	if uui.Name != "" {
		user.Name = &uui.Name
	}

	if err := update(tx, user); err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func (m *Manager) isValidURL(toTest string) bool {
	_, err := url.ParseRequestURI(toTest)
	if err != nil {
		return false
	}
	return true
}

func (m *Manager) isValidBase64(toTest string) bool {
	_, err := b64.RawStdEncoding.DecodeString(toTest)
	if err != nil {
		return false
	}
	return true
}
