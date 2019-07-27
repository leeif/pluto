package manage

import (
	"errors"

	"github.com/leeif/pluto/models"

	saltUtil "github.com/leeif/pluto/utils/salt"

	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/refresh"

	"github.com/jinzhu/gorm"
)

func LoginWithEmail(db *gorm.DB, login request.MailLogin) (string, *PlutoError) {
	if db == nil {
		return "", newPlutoError(ServerError, errors.New("DB connection is empty"))
	}

	tx := db.Begin()

	user := models.User{}
	if tx.Where("email = ?", login.Mail).First(&user).RecordNotFound() {
		return "", newPlutoError(ServerError, errors.New("email is not existed"))
	}

	salt := models.Salt{}
	if tx.Where("user_id = ?", user.ID).First(&salt).RecordNotFound() {
		return "", newPlutoError(ServerError, errors.New("Salt is not found"))
	}

	s, err := saltUtil.DecodeSalt(salt.Salt)

	if err != nil {
		return "", newPlutoError(ServerError, errors.New("Salt decoding is failed: "+err.Error()))
	}

	encodePassword, err := saltUtil.EncodePassword(login.Password, s)

	if err != nil {
		return "", newPlutoError(ServerError, errors.New("Password encoding is failed: "+err.Error()))
	}

	if user.Password != encodePassword {
		return "", newPlutoError(ServerError, errors.New("Password is invalid"))
	}

	// insert deviceID and appID into device table
	device := models.Device{}

	if tx.Where("device_id = ? and app_id = ?").First(&device).RecordNotFound() {
		device.DeviceID = login.DeviceID
		device.AppID = login.AppID
		if err := create(tx, &device); err != nil {
			return "", err
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
			return "", err
		}
	} else {
		rt.RefreshToken = refreshToken
		if err := update(tx, &device); err != nil {
			return "", err
		}
	}

	// generate jwt token
	jwtToken, err := jwt.GenerateUserJWT(jwt.Head{Alg: jwt.ALGRAS},
		jwt.UserPayload{UserID: user.ID, DeviceID: device.DeviceID, AppID: device.AppID})

	if err != nil {
		return "", &PlutoError{
			Type: ServerError,
			Err:  errors.New("jwt token generate failed: " + err.Error()),
		}
	}

	return jwtToken, nil
}

func RegisterWithEmail(db *gorm.DB, register request.MailRegister) *PlutoError {
	if db == nil {
		return &PlutoError{
			Type: ServerError,
			Err:  errors.New("DB connection is empty"),
		}
	}

	user := models.User{}
	if !db.Where("email = ?", register.Mail).First(&user).RecordNotFound() {
		return newPlutoError(ReqError, errors.New("Email is already exists"))
	}

	salt := models.Salt{}
	salt.Salt = saltUtil.RandomSalt(register.Mail)

	encodedPassword, err := saltUtil.EncodePassword(register.Password, salt.Salt)
	if err != nil {
		return newPlutoError(ServerError, errors.New("salt encoding is failed"))
	}

	user.Password = encodedPassword

	return nil
}

func ResetPassword() error {
	return nil
}
