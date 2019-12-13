package manage

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/leeif/pluto/utils/avatar"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/modelexts"
	"github.com/leeif/pluto/models"

	saltUtil "github.com/leeif/pluto/utils/salt"

	b64 "encoding/base64"

	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/utils/jwt"
)

func (m *Manager) ResetPasswordMail(rpm request.ResetPasswordMail) (*models.User, *perror.PlutoError) {

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(rpm.Mail))
	user, err := models.Users(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).One(m.db)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.MailNotExist
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	return user, nil
}

func (m *Manager) ResetPasswordPage(token string) *perror.PlutoError {

	jwtToken, perr := jwt.VerifyB64JWT(token)
	// token verify failed
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

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(prp.Mail))
	user, err := models.Users(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).One(m.db)
	if err != nil && err == sql.ErrNoRows {
		return perror.ServerError.Wrapper(errors.New("mail not found"))
	} else if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	// user is updated after password reset token is created
	if user.UpdatedAt.Valid && user.UpdatedAt.Time.Unix() > prp.Create {
		return perror.InvalidJWTToekn
	}

	return nil
}

func (m *Manager) ResetPassword(token string, rp request.ResetPasswordWeb) *perror.PlutoError {

	jwtToken, perr := jwt.VerifyB64JWT(token)
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

	tx, err := m.db.Begin()

	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(prp.Mail))
	user, err := models.Users(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).One(tx)
	if err != nil && err == sql.ErrNoRows {
		return perror.ServerError.Wrapper(errors.New("mail not found"))
	} else if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	// user is updated after password reset token is created
	if user.UpdatedAt.Valid && user.UpdatedAt.Time.Unix() > prp.Create {
		return perror.InvalidJWTToekn
	}

	salt, err := models.Salts(qm.Where("user_id = ?", user.ID)).One(tx)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	salt.Salt = saltUtil.RandomSalt(prp.Mail)

	if _, err := salt.Update(tx, boil.Whitelist("salt")); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	encodedPassword, perr := saltUtil.EncodePassword(rp.Password, salt.Salt)
	if perr != nil {
		return perror.ServerError.Wrapper(errors.New("Salt encoding is failed"))
	}

	user.Password.SetValid(encodedPassword)
	if _, err := user.Update(tx, boil.Infer()); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	tx.Commit()

	return nil
}

func (m *Manager) UserInfo(accessPayload *jwt.AccessPayload) (*modelexts.User, *perror.PlutoError) {

	user, err := models.Users(qm.Where("id = ?", accessPayload.UserID)).One(m.db)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(errors.New("user not found id: " + string(accessPayload.UserID)))
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	role, perr := getUserRole(accessPayload.UserID, accessPayload.AppID, m.db)
	if perr != nil {
		return nil, perr
	}

	userExt := &modelexts.User{
		User: user,
	}

	if role != nil {
		userExt.Roles = []string{role.Name}
	}

	return userExt, nil
}

func (m *Manager) RefreshAccessToken(rat request.RefreshAccessToken) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)

	tx, err := m.db.Begin()

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	rt, err := models.RefreshTokens(qm.Where("user_id = ? and refresh_token = ?", rat.UseID, rat.RefreshToken)).One(tx)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.InvalidRefreshToken
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	da, err := models.DeviceApps(qm.Where("id = ?", rt.DeviceAppID)).One(tx)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.InvalidRefreshToken
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	if rt.UserID != rat.UseID || da.DeviceID != rat.DeviceID || da.AppID != rat.AppID {
		return nil, perror.InvalidRefreshToken
	}

	user, err := models.Users(qm.Where("id = ?", rat.UseID)).One(tx)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(fmt.Errorf("UserID not found: %d", rat.UseID))
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	scopes := strings.Split(rt.Scopes.String, ",")

	// generate jwt token
	up := jwt.NewAccessPayload(rat.UseID, scopes, rat.AppID, user.LoginType, m.config.JWT.AccessTokenExpire)
	token, perr := jwt.GenerateRSAJWT(up)

	if perr != nil {
		return nil, perr.Wrapper(errors.New("JWT token generate failed"))
	}

	res["jwt"] = token.String()

	tx.Commit()
	return res, nil
}

func (m *Manager) UpdateUserInfo(accessPayload *jwt.AccessPayload, uui request.UpdateUserInfo) *perror.PlutoError {

	// only user using mail login can be update
	if accessPayload.LoginType != MAILLOGIN {
		return perror.InvalidJWTToekn
	}

	tx, err := m.db.Begin()

	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	user, err := models.Users(qm.Where("id = ?", accessPayload.UserID)).One(tx)
	if err != nil && err == sql.ErrNoRows {
		return perror.ServerError.Wrapper(errors.New("user not found id: " + string(accessPayload.UserID)))
	} else if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	if uui.Avatar != "" && m.isValidURL(uui.Avatar) {
		user.Avatar.SetValid(uui.Avatar)
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
		user.Avatar.SetValid(url)
	} else if uui.Avatar != "" {
		return perror.InvalidAvatarFormat
	}

	if uui.Name != "" {
		user.Name = uui.Name
	}

	if _, err := user.Update(tx, boil.Infer()); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	tx.Commit()

	return nil
}

func (m *Manager) isValidURL(toTest string) bool {
	u, err := url.Parse(toTest)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (m *Manager) isValidBase64(toTest string) bool {
	_, err := b64.RawStdEncoding.DecodeString(toTest)
	if err != nil {
		return false
	}
	return true
}
