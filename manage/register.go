package manage

import (
	"encoding/json"
	"errors"
	"time"

	b64 "encoding/base64"

	"github.com/jinzhu/gorm"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/models"
	"github.com/leeif/pluto/utils/avatar"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/mail"
	saltUtil "github.com/leeif/pluto/utils/salt"
)

func (m *Manger) RegisterWithEmail(register request.MailRegister) (uint, *perror.PlutoError) {

	tx := m.db.Begin()
	defer func() {
		tx.Rollback()
	}()

	user := models.User{}
	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(register.Mail))
	if !tx.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken).First(&user).RecordNotFound() {
		return 0, perror.MailIsAlreadyRegister
	}

	salt := models.Salt{}
	salt.Salt = saltUtil.RandomSalt(identifyToken)

	encodedPassword, err := saltUtil.EncodePassword(register.Password, salt.Salt)
	if err != nil {
		return 0, perror.ServerError.Wrapper(errors.New("Salt encoding is failed"))
	}

	user.Mail = &register.Mail
	user.Name = &register.Name
	user.IdentifyToken = identifyToken
	user.LoginType = MAILLOGIN
	user.Password = &encodedPassword

	if m.config.Server.SkipRegisterVerifyMail {
		user.Verified = true
	}

	// get a random avatar
	a := avatar.NewAvatar(m.config)
	avatarURL, err := a.GetRandomAvatar()
	user.Avatar = avatarURL
	if err != nil {
		m.logger.Warn(err.LogError)
	}

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

func (m *Manger) RegisterVerifyMail(db *gorm.DB, rvm request.RegisterVerifyMail, baseURL string) *perror.PlutoError {
	if db == nil {
		return perror.ServerError.Wrapper(errors.New("DB connection is empty"))
	}

	user := models.User{}
	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(rvm.Mail))
	if db.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken).First(&user).RecordNotFound() {
		return perror.MailIsNotExsit
	}

	if user.Verified == true {
		return perror.MailAlreadyVerified
	}

	ml, err := mail.NewMail(m.config)
	if err != nil {
		return err
	}
	if err := ml.SendRegisterVerify(user.ID, *user.Mail, baseURL); err != nil {
		return err
	}

	return nil
}

func (m *Manger) RegisterVerify(token string) *perror.PlutoError {

	jwtToken, perr := jwt.VerifyB64JWT(token)
	if perr != nil {
		return perr
	}

	verifyPayload := jwt.RegisterVerifyPayload{}
	if err := json.Unmarshal(jwtToken.Payload, &verifyPayload); err != nil {
		return perror.ServerError.Wrapper(errors.New("parse user payload failed")).Wrapper(err)
	}

	if verifyPayload.Type != jwt.REGISTERVERIFY {
		return perror.InvalidJWTToekn
	}

	// expire
	if time.Now().Unix() > verifyPayload.Expire {
		return perror.InvalidJWTToekn
	}

	tx := m.db.Begin()
	defer func() {
		tx.Rollback()
	}()

	user := models.User{}
	if tx.Where("id = ?", verifyPayload.UserID).First(&user).RecordNotFound() {
		return perror.ServerError.Wrapper(errors.New("user not found"))
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
