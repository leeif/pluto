package manage

import (
	"crypto"
	"errors"
	"strconv"
	"strings"

	b64 "encoding/base64"
	"encoding/json"

	"github.com/jinzhu/gorm"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/models"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/rsa"
)

func RefreshAccessToken(db *gorm.DB, rat request.RefreshAccessToken) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)

	if db == nil {
		return nil, perror.NewServerError(errors.New("DB connection is empty"))
	}

	if !rat.Validation() {
		return nil, perror.BadRequest
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
		return nil, perror.NewServerError(errors.New("Refresh token is valid but user id is not existed: " + strconv.Itoa(int(rt.UserID))))
	}

	// generate jwt token
	jwtToken, err := jwt.GenerateUserJWT(jwt.Head{Alg: jwt.ALGRAS},
		jwt.UserPayload{UserID: user.ID, DeviceID: rat.DeviceID, AppID: rat.AppID}, 60*60)

	if err != nil {
		return nil, perror.NewServerError(errors.New("JWT token generate failed: " + err.Error()))
	}

	res["jwt"] = jwtToken

	tx.Commit()
	return res, nil
}

func VerifyMail(db *gorm.DB, token string) *perror.PlutoError {

	b, err := b64.StdEncoding.DecodeString(token)
	if err != nil {
		return perror.NewServerError(errors.New("base 64 decode failed: " + err.Error()))
	}
	jwtToken := strings.Split(string(b), ".")
	header, err := b64.StdEncoding.DecodeString(jwtToken[0])
	if err != nil {
		return perror.NewServerError(errors.New("base 64 decode failed: " + err.Error()))
	}
	payload, err := b64.StdEncoding.DecodeString(jwtToken[1])
	if err != nil {
		return perror.NewServerError(errors.New("base 64 decode failed: " + err.Error()))
	}
	signed, err := b64.StdEncoding.DecodeString(jwtToken[2])
	if err != nil {
		return perror.NewServerError(errors.New("base 64 decode failed: " + err.Error()))
	}
	if err := rsa.VerifySignWithPublicKey(append(header, payload...), signed, crypto.SHA256); err != nil {
		return perror.InvalidMailVerifyToekn
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
