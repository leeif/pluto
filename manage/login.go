package manage

import (
	b64 "encoding/base64"
	"errors"
	"net/http"

	gjwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/models"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/refresh"
	saltUtil "github.com/leeif/pluto/utils/salt"
	"google.golang.org/api/oauth2/v2"
)

const (
	MAILLOGIN   = "mail"
	GOOGLELOGIN = "google"
)

func LoginWithEmail(db *gorm.DB, login request.MailLogin) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)
	if db == nil {
		return nil, perror.ServerError.Wrapper(errors.New("DB connection is empty"))
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()

	user := models.User{}
	identifyToken := b64.StdEncoding.EncodeToString([]byte(login.Mail))
	if tx.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken).First(&user).RecordNotFound() {
		return nil, perror.MailIsNotExsit
	}

	if user.Verified == false {
		return nil, perror.MailIsNotVerified
	}

	salt := models.Salt{}
	if tx.Where("user_id = ?", user.ID).First(&salt).RecordNotFound() {
		return nil, perror.ServerError.Wrapper(errors.New("Salt is not found"))
	}

	encodePassword, err := saltUtil.EncodePassword(login.Password, salt.Salt)

	if err != nil {
		return nil, err.Wrapper(errors.New("Password encoding is failed"))
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
		return nil, err.Wrapper(errors.New("JWT token generate failed"))
	}

	res["jwt"] = jwtToken
	res["refresh_token"] = rt.RefreshToken

	tx.Commit()

	return res, nil
}

func LoginWithGoogle(db *gorm.DB, login request.GoogleLogin) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)

	info, err := verifyGoogleIdToken(login.IDToken)
	if err != nil {
		return nil, err
	}

	tx := db.Begin()
	defer func() {
		tx.Rollback()
	}()
	user := models.User{}
	user.IdentifyToken = info.Sub
	user.LoginType = GOOGLELOGIN
	user.Avatar = info.Picture
	user.Name = &info.Name
	user.Mail = &info.Email
	if tx.Where("login_type = ? and identify_token = ?", user.LoginType, user.IdentifyToken).First(&user).RecordNotFound() {
		if err := create(tx, &user); err != nil {
			return nil, err
		}
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
		return nil, err.Wrapper(errors.New("JWT token generate failed"))
	}

	res["jwt"] = jwtToken
	res["refresh_token"] = rt.RefreshToken

	tx.Commit()
	return res, nil
}

// googleIDTokenInfo struct
type googleIDTokenInfo struct {
	Iss string `json:"iss"`
	// userId
	Sub string `json:"sub"`
	Azp string `json:"azp"`
	// clientId
	Aud string `json:"aud"`
	Iat int64  `json:"iat"`
	// expired time
	Exp int64 `json:"exp"`

	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Local         string `json:"locale"`
	gjwt.StandardClaims
}

func verifyGoogleIdToken(idToken string) (*googleIDTokenInfo, *perror.PlutoError) {
	var httpClient = &http.Client{}
	oauth2Service, err := oauth2.New(httpClient)
	tokenInfoCall := oauth2Service.Tokeninfo()
	tokenInfoCall.IdToken(idToken)
	tokenInfo, err := tokenInfoCall.Do()
	if err != nil {
		return nil, perror.InvalidIDToken.Wrapper(err)
	}
	if tokenInfo.Audience != "" {
		return nil, perror.InvalidIDToken
	}
	parser := gjwt.Parser{}
	token, _, err := parser.ParseUnverified(idToken, &googleIDTokenInfo{})
	if err != nil {
		return nil, perror.InvalidIDToken.Wrapper(err)
	}
	if info, ok := token.Claims.(*googleIDTokenInfo); ok {
		return info, nil
	}
	return nil, perror.InvalidIDToken
}
