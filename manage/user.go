package manage

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/models"

	saltUtil "github.com/leeif/pluto/utils/salt"

	b64 "encoding/base64"

	"github.com/jinzhu/gorm"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/mail"
)

func ResetPasswordMail(db *gorm.DB, rpm request.ResetPasswordMail) *perror.PlutoError {
	if db == nil {
		return perror.ServerError.Wrapper(errors.New("DB connection is empty"))
	}

	user := models.User{}
	identifyToken := b64.StdEncoding.EncodeToString([]byte(rpm.Mail))
	if db.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken).First(&user).RecordNotFound() {
		return perror.MailIsNotExsit
	}

	if err := mail.SendResetPassword(*user.Mail); err != nil {
		return err
	}

	return nil
}

func ResetPasswordPage(db *gorm.DB, token string) *perror.PlutoError {
	if db == nil {
		return perror.ServerError.Wrapper(errors.New("DB connection is empty"))
	}

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
	if db.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken).First(&user).RecordNotFound() {
		return perror.ServerError.Wrapper(errors.New("mail not found"))
	}

	// user is updated after password reset token is created
	if user.UpdatedAt.Unix() > prp.Create {
		return perror.InvalidJWTToekn
	}

	return nil
}

func ResetPassword(db *gorm.DB, rp request.ResetPassword) *perror.PlutoError {
	if db == nil {
		return perror.ServerError.Wrapper(errors.New("DB connection is empty"))
	}

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

	tx := db.Begin()
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

	tx.Commit()

	return nil
}

func UserInfo(db *gorm.DB, token string) (*models.User, *perror.PlutoError) {
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
	if db.Where("id = ?", userPayload.UserID).First(&user).RecordNotFound() {
		return nil, perror.ServerError.Wrapper(errors.New("user not found id: " + string(userPayload.UserID)))
	}

	return &user, nil
}

func RefreshAccessToken(db *gorm.DB, rat request.RefreshAccessToken) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)

	if db == nil {
		return nil, perror.ServerError.Wrapper(errors.New("DB connection is empty"))
	}

	tx := db.Begin()

	defer func() {
		tx.Rollback()
	}()

	rt := models.RefreshToken{}
	if tx.Where("refresh_token = ?", rat.RefreshToken).First(&rt).RecordNotFound() {
		return nil, perror.InvalidRefreshToken
	}

	if rt.UserID != rat.UseID || rt.DeviceID != rat.DeviceID || rt.AppID != rat.AppID {
		return nil, perror.InvalidRefreshToken
	}

	user := models.User{}
	if tx.Where("id = ?", rt.UserID).First(&user).RecordNotFound() {
		return nil, perror.ServerError.Wrapper(errors.New("Refresh token is valid but user id is not existed: " + strconv.Itoa(int(rt.UserID))))
	}

	// generate jwt token
	jwtToken, err := jwt.GenerateJWT(jwt.Head{Type: jwt.ACCESS},
		&jwt.UserPayload{UserID: user.ID, DeviceID: rat.DeviceID, AppID: rat.AppID}, 60*60)

	if err != nil {
		return nil, err.Wrapper(errors.New("JWT token generate failed"))
	}

	res["jwt"] = jwtToken

	tx.Commit()
	return res, nil
}
