package manage

import (
	"errors"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/leeif/pluto/datatype"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/models"
	"github.com/leeif/pluto/utils/jwt"
)

func RefreshAccessToken(db *gorm.DB, rat request.RefreshAccessToken) (map[string]string, *datatype.PlutoError) {
	res := make(map[string]string)

	if db == nil {
		return nil, datatype.NewPlutoError(datatype.ServerError,
			errors.New("DB connection is empty"))
	}

	if !rat.Validation() {
		return nil, datatype.NewPlutoError(datatype.ReqError,
			errors.New("Request parameters are not enough"))
	}

	tx := db.Begin()
	rt := models.RefreshToken{}
	if tx.Where("refresh_token = ?", rat.RefreshToken).First(&rt).RecordNotFound() {
		return nil, datatype.NewPlutoError(datatype.ReqError,
			errors.New("Invalid refresh token"))
	}

	if rt.UserID != rat.UseID || rt.DeviceID != rat.DeviceID || rt.AppID != rat.AppID {
		return nil, datatype.NewPlutoError(datatype.ReqError,
			errors.New("Invalid refresh token"))
	}

	user := models.User{}
	if tx.Where("id = ?", rt.UserID).First(&user).RecordNotFound() {
		return nil, datatype.NewPlutoError(datatype.ServerError,
			errors.New("Refresh token is valid but user id is not existed: "+strconv.Itoa(int(rt.UserID))))
	}

	// generate jwt token
	jwtToken, err := jwt.GenerateUserJWT(jwt.Head{Alg: jwt.ALGRAS},
		jwt.UserPayload{UserID: user.ID, DeviceID: rat.DeviceID, AppID: rat.AppID})

	if err != nil {
		return nil, datatype.NewPlutoError(datatype.ServerError,
			errors.New("JWT token generate failed: "+err.Error()))
	}

	res["jwt"] = jwtToken

	tx.Commit()
	return res, nil
}
