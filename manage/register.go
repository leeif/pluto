package manage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/volatiletech/sqlboiler/queries/qm"

	b64 "encoding/base64"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/models"
	"github.com/leeif/pluto/utils/avatar"
	"github.com/leeif/pluto/utils/jwt"
	saltUtil "github.com/leeif/pluto/utils/salt"
)

func (m *Manager) RegisterWithEmail(register request.MailRegister) (*models.User, *perror.PlutoError) {

	tx, err := m.db.Begin()
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}
	defer func() {
		tx.Rollback()
	}()

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(register.Mail))
	user, err := models.Users(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	if user != nil {
		return user, perror.MailIsAlreadyRegister
	}

	salt := saltUtil.RandomSalt(identifyToken)

	encodedPassword, perr := saltUtil.EncodePassword(register.Password, salt)
	if perr != nil {
		return nil, perr
	}
	user = &models.User{}
	user.Mail.SetValid(register.Mail)
	user.Name = register.Name
	user.IdentifyToken = identifyToken
	user.LoginType = MAILLOGIN
	user.Password.SetValid(encodedPassword)

	if m.config.Server.SkipRegisterVerifyMail {
		user.Verified.SetValid(true)
	}

	// get a random avatar
	ag := avatar.AvatarGen{}
	avatarReader, perr := ag.GenFromGravatar()
	if perr != nil {
		return nil, perr
	}

	as := avatar.NewAvatarSaver(m.config)
	remoteURL, perr := as.SaveAvatarImageInOSS(avatarReader)
	if perr != nil {
		user.Avatar.SetValid(avatarReader.OriginURL)
		m.logger.Error(perr.LogError)
	} else {
		user.Avatar.SetValid(remoteURL)
	}

	if err := user.Insert(tx, boil.Infer()); err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	saltModel := models.Salt{}
	saltModel.Salt = salt
	saltModel.UserID = user.ID
	if err := saltModel.Insert(tx, boil.Infer()); err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	tx.Commit()

	return user, nil
}

func (m *Manager) RegisterVerifyMail(rvm request.RegisterVerifyMail) (*models.User, *perror.PlutoError) {

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(rvm.Mail))
	user, err := models.Users(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).One(m.db)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.MailNotExist
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	if user.Verified.Bool == true {
		return nil, perror.MailAlreadyVerified
	}

	return user, nil
}

func (m *Manager) RegisterVerify(token string) *perror.PlutoError {

	jwtToken, perr := jwt.VerifyB64JWT(token)
	if perr != nil {
		return perr
	}

	verifyPayload := jwt.RegisterVerifyPayload{}
	if err := json.Unmarshal(jwtToken.Payload, &verifyPayload); err != nil {
		return perror.ServerError.Wrapper(errors.New("parse user payload failed")).Wrapper(err)
	}

	if verifyPayload.Type != jwt.REGISTERVERIFY {
		return perror.InvalidJWTToken
	}

	// expire
	if time.Now().Unix() > verifyPayload.Expire {
		return perror.JWTTokenExpired
	}

	tx, err := m.db.Begin()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	user, err := models.Users(qm.Where("id = ?", verifyPayload.UserID)).One(tx)
	if err != nil && err == sql.ErrNoRows {
		return perror.ServerError.Wrapper(errors.New("user not found"))
	} else if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	if user.Verified.Bool == true {
		return perror.MailAlreadyVerified
	}

	user.Verified.SetValid(true)

	user.Update(tx, boil.Whitelist("verified"))

	tx.Commit()

	return nil
}
