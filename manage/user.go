package manage

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/alecthomas/template"
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
		if r := recover(); r != nil {
			tx.Rollback()
		}
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
	jwtToken, err := jwt.GenerateUserJWT(jwt.Head{Alg: jwt.ALGRAS},
		jwt.UserPayload{UserID: user.ID, DeviceID: device.DeviceID, AppID: device.AppID}, 60*60)

	if err != nil {
		return nil, perror.NewServerError(errors.New("JWT token generate failed: " + err.Error()))
	}

	res["jwt"] = jwtToken
	res["refresh_token"] = rt.RefreshToken

	tx.Commit()

	return res, nil
}

func RegisterWithEmail(db *gorm.DB, register request.MailRegister) *perror.PlutoError {
	if db == nil {
		return perror.NewServerError(errors.New("DB connection is empty"))
	}

	if !register.Validation() {
		return perror.BadRequest
	}

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user := models.User{}
	if !tx.Where("mail = ?", register.Mail).First(&user).RecordNotFound() {
		return perror.MailIsAlreadyRegister
	}

	salt := models.Salt{}
	salt.Salt = saltUtil.RandomSalt(register.Mail)

	encodedPassword, err := saltUtil.EncodePassword(register.Password, salt.Salt)
	if err != nil {
		return perror.NewServerError(errors.New("Salt encoding is failed"))
	}

	user.Mail = &register.Mail
	user.Name = &register.Name
	user.Password = &encodedPassword

	if err := create(tx, &user); err != nil {
		return err
	}

	salt.UserID = user.ID

	if err := create(tx, &salt); err != nil {
		return err
	}

	token, err := jwt.GenerateUserJWT(jwt.Head{Alg: jwt.ALGRAS}, jwt.UserPayload{UserID: user.ID}, 10*60)
	if err != nil {
		return perror.NewServerError(errors.New("JWT token generate failed: " + err.Error()))
	}

	// send verify mail (must)
	if m := mail.NewMail(); m != nil {
		dir, _ := os.Getwd()
		t := template.Must(template.ParseFiles(path.Join(dir, "views", "mail_verify.html")))
		var buffer bytes.Buffer
		type Data struct {
			Token string
		}
		t.Execute(&buffer, Data{Token: token})
		fmt.Println(buffer.String())
		if err := m.Send(register.Mail, "[MuShare]Mail Verification", "test"); err != nil {
			tx.Rollback()
			return perror.NewServerError(errors.New("Mail sending failed: " + err.Error()))
		}
	} else {
		tx.Rollback()
		return perror.NewServerError(errors.New("Mail sender is not defined"))
	}

	tx.Commit()

	return nil
}

func ResetPassword() error {
	return nil
}
