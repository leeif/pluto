package manage

import (
	"encoding/json"
	"errors"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/models"

	saltUtil "github.com/leeif/pluto/utils/salt"

	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/mail"
	"github.com/leeif/pluto/utils/refresh"

	"github.com/jinzhu/gorm"
)

func LoginWithEmail(db *gorm.DB, login request.MailLogin) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)
	if db == nil {
		return nil, perror.NewServerError(errors.New("DB connection is empty"))
	}

	if !login.Validation() {
		return nil, perror.BadRequest
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	user := models.User{}
	if tx.Where("mail = ?", login.Mail).First(&user).RecordNotFound() {
		return nil, perror.MailIsNotExsit
	}

	if user.Verified == false {
		return nil, perror.MailIsNotVerified
	}

	salt := models.Salt{}
	if tx.Where("user_id = ?", user.ID).First(&salt).RecordNotFound() {
		return nil, perror.NewServerError(errors.New("Salt is not found"))
	}

	encodePassword, err := saltUtil.EncodePassword(login.Password, salt.Salt)

	if err != nil {
		return nil, perror.NewServerError(errors.New("Password encoding is failed: " + err.Error()))
	}

	if *user.Password != encodePassword {
		return nil, perror.InvalidPassword
	}

	// insert deviceID and appID into device table
	device := models.Device{}

	if tx.Where("device_id = ? and app_id = ?", login.DeviceID, login.AppID).First(&device).RecordNotFound() {
		device.DeviceID = login.DeviceID
		device.AppID = login.AppID
		if err := create(tx, &device); err != nil {
			return nil, err
		}
	}

	// refresh token
	rt := models.RefreshToken{}
	refreshToken := refresh.GenerateRefreshToken(string(user.ID) + device.DeviceID + device.AppID)
	rt.UserID = user.ID
	rt.DeviceID = device.DeviceID
	rt.AppID = device.AppID
	if tx.Where("device_id = ? and app_id = ? and user_id = ?", device.DeviceID, device.AppID, user.ID).First(&rt).RecordNotFound() {
		rt.RefreshToken = refreshToken
		if err := create(tx, &rt); err != nil {
			return nil, err
		}
	} else {
		rt.RefreshToken = refreshToken
		if err := update(tx, &rt); err != nil {
			return nil, err
		}
	}

	// generate jwt token
	jwtToken, err := jwt.GenerateJWT(jwt.Head{Type: jwt.ACCESS},
		&jwt.UserPayload{UserID: user.ID, DeviceID: device.DeviceID, AppID: device.AppID}, 60*60)

	if err != nil {
		return nil, perror.NewServerError(errors.New("JWT token generate failed: " + err.Error()))
	}

	res["jwt"] = jwtToken
	res["refresh_token"] = rt.RefreshToken

	tx.Commit()

	return res, nil
}

func RegisterWithEmail(db *gorm.DB, register request.MailRegister) (uint, *perror.PlutoError) {
	if db == nil {
		return 0, perror.NewServerError(errors.New("DB connection is empty"))
	}

	if !register.Validation() {
		return 0, perror.BadRequest
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	user := models.User{}
	if !tx.Where("mail = ?", register.Mail).First(&user).RecordNotFound() {
		return 0, perror.MailIsAlreadyRegister
	}

	salt := models.Salt{}
	salt.Salt = saltUtil.RandomSalt(register.Mail)

	encodedPassword, err := saltUtil.EncodePassword(register.Password, salt.Salt)
	if err != nil {
		return 0, perror.NewServerError(errors.New("Salt encoding is failed"))
	}

	user.Mail = &register.Mail
	user.Name = &register.Name
	user.Password = &encodedPassword

	if err := create(tx, &user); err != nil {
		return 0, err
	}

	salt.UserID = user.ID

	if err := create(tx, &salt); err != nil {
		return 0, err
	}

	tx.Commit()

	return user.ID, nil
}

func RegisterVerifyMail(db *gorm.DB, rvm request.RegisterVerifyMail) *perror.PlutoError {
	if db == nil {
		return perror.NewServerError(errors.New("DB connection is empty"))
	}

	if !rvm.Validation() {
		return perror.BadRequest
	}

	user := models.User{}
	if db.Where("mail = ?", rvm.Mail).First(&user).RecordNotFound() {
		return perror.MailIsNotExsit
	}

	if user.Verified == true {
		return perror.MailAlreadyVerified
	}

	if err := mail.SendRegisterVerify(user.ID, *user.Mail); err != nil {
		return err
	}

	return nil
}

func RegisterVerify(db *gorm.DB, token string) *perror.PlutoError {

	header, payload, perr := jwt.VerifyB64JWT(token)
	if perr != nil {
		return perr
	}

	head := jwt.Head{}
	if err := json.Unmarshal(header, &head); err != nil {
		return perror.NewServerError(errors.New("parse password reset payload failed: " + err.Error()))
	}

	if head.Type != jwt.REGISTERVERIFY {
		return perror.InvalidJWTToekn
	}

	userPayload := jwt.UserPayload{}
	if err := json.Unmarshal(payload, &userPayload); err != nil {
		return perror.NewServerError(errors.New("parse user payload failed: " + err.Error()))
	}

	// currently no expired
	// if time.Now().Unix() > userPayload.Expire {

	// }

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	user := models.User{}
	if tx.Where("id = ?", userPayload.UserID).First(&user).RecordNotFound() {
		return perror.NewServerError(errors.New("user not found"))
	}

	if user.Verified == true {
		return perror.MailAlreadyVerified
	}

	user.Verified = true

	if err := update(tx, &user); err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func ResetPasswordMail(db *gorm.DB, rpm request.ResetPasswordMail) *perror.PlutoError {
	if db == nil {
		return perror.NewServerError(errors.New("DB connection is empty"))
	}

	if !rpm.Validation() {
		return perror.BadRequest
	}

	user := models.User{}
	if db.Where("mail = ?", rpm.Mail).First(&user).RecordNotFound() {
		return perror.MailIsNotExsit
	}

	if err := mail.SendResetPassword(*user.Mail); err != nil {
		return err
	}

	return nil
}

func ResetPasswordPage(db *gorm.DB, token string) *perror.PlutoError {
	if db == nil {
		return perror.NewServerError(errors.New("DB connection is empty"))
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
	if err := json.Unmarshal(payload, &prp); err != nil {
		return perror.NewServerError(errors.New("parse password reset payload failed: " + err.Error()))
	}

	user := models.User{}
	if db.Where("mail = ?", prp.Mail).First(&user).RecordNotFound() {
		return perror.NewServerError(errors.New("mail not found"))
	}

	// user is updated after password reset token is created
	if user.UpdatedAt.Unix() > prp.Create {
		return perror.InvalidJWTToekn
	}

	return nil
}

func ResetPassword(db *gorm.DB, rp request.ResetPassword) *perror.PlutoError {
	if db == nil {
		return perror.NewServerError(errors.New("DB connection is empty"))
	}

	if !rp.Validation() {
		return perror.BadRequest
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

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	user := models.User{}
	if tx.Where("mail = ?", prp.Mail).First(&user).RecordNotFound() {
		return perror.NewServerError(errors.New("mail not found"))
	}

	// user is updated after password reset token is created
	if user.UpdatedAt.Unix() > prp.Create {
		return perror.InvalidJWTToekn
	}

	salt := models.Salt{}
	if tx.Where("user_id = ?", user.ID).First(&salt).RecordNotFound() {
		return perror.NewServerError(errors.New("Salt is not found"))
	}

	salt.Salt = saltUtil.RandomSalt(prp.Mail)

	if err := update(tx, &salt); err != nil {
		return err
	}

	encodedPassword, err := saltUtil.EncodePassword(rp.Password, salt.Salt)
	if err != nil {
		return perror.NewServerError(errors.New("Salt encoding is failed"))
	}

	user.Password = &encodedPassword
	if err := update(tx, &user); err != nil {
		return err
	}

	tx.Commit()

	return nil
}
