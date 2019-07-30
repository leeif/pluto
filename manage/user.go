package manage

import (
	"errors"

	"github.com/leeif/pluto/datatype"
	"github.com/leeif/pluto/models"

	saltUtil "github.com/leeif/pluto/utils/salt"

	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/refresh"

	"github.com/jinzhu/gorm"
)

func LoginWithEmail(db *gorm.DB, login request.MailLogin) (map[string]string, *datatype.PlutoError) {
	res := make(map[string]string)
	if db == nil {
		return nil, datatype.NewPlutoError(datatype.ServerError,
			errors.New("DB connection is empty"))
	}

	if !login.Validation() {
		return nil, datatype.NewPlutoError(datatype.ReqError,
			errors.New("Request parameters are not enough"))
	}

	tx := db.Begin()

	user := models.User{}
	if tx.Where("mail = ?", login.Mail).First(&user).RecordNotFound() {
		return nil, datatype.NewPlutoError(datatype.ReqError,
			errors.New("Mail is not existed"))
	}

	salt := models.Salt{}
	if tx.Where("user_id = ?", user.ID).First(&salt).RecordNotFound() {
		return nil, datatype.NewPlutoError(datatype.ServerError,
			errors.New("Salt is not found"))
	}

	encodePassword, err := saltUtil.EncodePassword(login.Password, salt.Salt)

	if err != nil {
		return nil, datatype.NewPlutoError(datatype.ServerError,
			errors.New("Password encoding is failed: "+err.Error()))
	}

	if user.Password != encodePassword {
		return nil, datatype.NewPlutoError(datatype.ReqError,
			errors.New("Password is invalid"))
	}

	// insert deviceID and appID into device table
	device := models.Device{}

	if tx.Where("device_id = ? and app_id = ?").First(&device).RecordNotFound() {
		device.DeviceID = login.DeviceID
		device.AppID = login.AppID
		if err := create(tx, &device); err != nil {
			return nil, err
		}
	}

	// refresh token
	rt := models.RefreshToken{}
	refreshToken := refresh.GenerateRefreshToken(string(user.ID) + device.DeviceID + device.AppID)
	rt.DeviceID = device.DeviceID
	rt.AppID = device.AppID
	if tx.Where("device_id = ? and app_id = ?", device.DeviceID, device.AppID).First(&rt).RecordNotFound() {
		rt.RefreshToken = refreshToken
		if err := create(tx, &rt); err != nil {
			return nil, err
		}
	} else {
		rt.RefreshToken = refreshToken
		if err := update(tx, &device); err != nil {
			return nil, err
		}
	}

	// generate jwt token
	jwtToken, err := jwt.GenerateUserJWT(jwt.Head{Alg: jwt.ALGRAS},
		jwt.UserPayload{UserID: user.ID, DeviceID: device.DeviceID, AppID: device.AppID})

	if err != nil {
		return nil, datatype.NewPlutoError(datatype.ServerError,
			errors.New("JWT token generate failed: "+err.Error()))
	}

	res["jwt"] = jwtToken
	res["refresh_token"] = rt.RefreshToken

	return res, nil
}

func RegisterWithEmail(db *gorm.DB, register request.MailRegister) *datatype.PlutoError {
	if db == nil {
		return datatype.NewPlutoError(datatype.ServerError,
			errors.New("DB connection is empty"))
	}

	if !register.Validation() {
		return datatype.NewPlutoError(datatype.ReqError,
			errors.New("Request parameters are not enough"))
	}

	tx := db.Begin()

	user := models.User{}
	if !tx.Where("mail = ?", register.Mail).First(&user).RecordNotFound() {
		return datatype.NewPlutoError(datatype.ReqError,
			errors.New("Mail is already exists"))
	}

	salt := models.Salt{}
	salt.Salt = saltUtil.RandomSalt(register.Mail)

	encodedPassword, err := saltUtil.EncodePassword(register.Password, salt.Salt)
	if err != nil {
		return datatype.NewPlutoError(datatype.ServerError,
			errors.New("Salt encoding is failed"))
	}

	user.Mail = register.Mail
	user.Password = encodedPassword

	if err := create(tx, &user); err != nil {
		return err
	}

	salt.UserID = user.ID

	if err := create(tx, &salt); err != nil {
		return err
	}

	tx.Commit()

	return nil
}

func ResetPassword() error {
	return nil
}
