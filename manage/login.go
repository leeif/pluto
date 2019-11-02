package manage

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/leeif/pluto/config"

	gjwt "github.com/dgrijalva/jwt-go"
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
	WECHATLOGIN = "wechat"
)

func (m *Manger) EmailLogin(login request.MailLogin) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)

	tx := m.db.Begin()
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
	deviceAPP := models.DeviceAPP{}

	if tx.Where("device_id = ? and app_id = ?", login.DeviceID, login.AppID).First(&deviceAPP).RecordNotFound() {
		deviceAPP.DeviceID = login.DeviceID
		deviceAPP.AppID = login.AppID
		if err := create(tx, &deviceAPP); err != nil {
			return nil, err.Wrapper(errors.New("table device_apps"))
		}
	}

	// refresh token
	rt := models.RefreshToken{}
	refreshToken := refresh.GenerateRefreshToken(string(user.ID) + deviceAPP.DeviceID + deviceAPP.AppID)
	rt.UserID = user.ID
	rt.DeviceAPPID = deviceAPP.ID
	if tx.Where("device_app_id = ? and user_id = ?", deviceAPP.ID, user.ID).First(&rt).RecordNotFound() {
		rt.RefreshToken = refreshToken
		if err := create(tx, &rt); err != nil {
			return nil, err.Wrapper(errors.New("table refresh_tokens"))
		}
	} else {
		rt.RefreshToken = refreshToken
		if err := update(tx, &rt); err != nil {
			return nil, err
		}
	}

	// generate jwt token
	up := jwt.NewUserPayload(user.ID, deviceAPP.DeviceID, deviceAPP.AppID, m.config.JWT.AccessTokenExpire)
	token, err := jwt.GenerateRSAJWT(up)

	if err != nil {
		return nil, err.Wrapper(errors.New("JWT token generate failed"))
	}

	// add operation history
	if err := historyOperation(tx, OperationMailLogin, user.ID); err != nil {
		return nil, err
	}

	res["jwt"] = token.String()
	res["refresh_token"] = rt.RefreshToken

	tx.Commit()

	return res, nil
}

func (m *Manger) GoogleLoginMobile(login request.GoogleMobileLogin) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)

	info, err := verifyGoogleIdToken(login.IDToken)
	if err != nil {
		return nil, err
	}

	tx := m.db.Begin()
	defer func() {
		tx.Rollback()
	}()
	user := models.User{}
	user.IdentifyToken = info.Sub
	user.LoginType = GOOGLELOGIN
	user.Avatar = info.Picture
	user.Name = &info.Name
	user.Mail = &info.Email
	user.Verified = true
	if tx.Where("login_type = ? and identify_token = ?", user.LoginType, user.IdentifyToken).First(&user).RecordNotFound() {
		if err := create(tx, &user); err != nil {
			return nil, err
		}
	}

	// insert deviceID and appID into device table
	deviceAPP := models.DeviceAPP{}

	if tx.Where("device_id = ? and app_id = ?", login.DeviceID, login.AppID).First(&deviceAPP).RecordNotFound() {
		deviceAPP.DeviceID = login.DeviceID
		deviceAPP.AppID = login.AppID
		if err := create(tx, &deviceAPP); err != nil {
			return nil, err
		}
	}

	// refresh token
	rt := models.RefreshToken{}
	refreshToken := refresh.GenerateRefreshToken(string(user.ID) + deviceAPP.DeviceID + deviceAPP.AppID)
	rt.UserID = user.ID
	rt.DeviceAPPID = deviceAPP.ID
	if tx.Where("device_app_id = ? and user_id = ?", deviceAPP.ID, user.ID).First(&rt).RecordNotFound() {
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
	up := jwt.NewUserPayload(user.ID, deviceAPP.DeviceID, deviceAPP.AppID, m.config.JWT.AccessTokenExpire)
	token, err := jwt.GenerateRSAJWT(up)

	if err != nil {
		return nil, err.Wrapper(errors.New("JWT token generate failed"))
	}

	res["jwt"] = token.String()
	res["refresh_token"] = rt.RefreshToken

	// add operation history
	if err := historyOperation(tx, OperationGoogleLogin, user.ID); err != nil {
		return nil, err
	}

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
		return nil, perror.InvalidGoogleIDToken.Wrapper(err)
	}
	if tokenInfo.Audience == "" {
		return nil, perror.InvalidGoogleIDToken
	}
	parser := gjwt.Parser{}
	token, _, err := parser.ParseUnverified(idToken, &googleIDTokenInfo{})
	if err != nil {
		return nil, perror.InvalidGoogleIDToken.Wrapper(err)
	}
	if info, ok := token.Claims.(*googleIDTokenInfo); ok {
		return info, nil
	}
	return nil, perror.InvalidGoogleIDToken
}

func (m *Manger) WechatLoginMobile(login request.WechatMobileLogin) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)

	accessToken, openID, err := getWechatAccessToken(login.Code, m.config.WechatLogin)

	if err != nil {
		return nil, err
	}

	info, err := getWebchatUserInfo(accessToken, openID)
	if err != nil {
		return nil, err
	}

	tx := m.db.Begin()
	defer func() {
		tx.Rollback()
	}()
	user := models.User{}
	user.IdentifyToken = info.Unionid
	user.LoginType = WECHATLOGIN
	user.Avatar = info.HeadimgURL
	user.Name = &info.Nickname
	user.Verified = true
	if tx.Where("login_type = ? and identify_token = ?", user.LoginType, user.IdentifyToken).First(&user).RecordNotFound() {
		if err := create(tx, &user); err != nil {
			return nil, err
		}
	}

	// insert deviceID and appID into device table
	deviceAPP := models.DeviceAPP{}

	if tx.Where("device_id = ? and app_id = ?", login.DeviceID, login.AppID).First(&deviceAPP).RecordNotFound() {
		deviceAPP.DeviceID = login.DeviceID
		deviceAPP.AppID = login.AppID
		if err := create(tx, &deviceAPP); err != nil {
			return nil, err
		}
	}

	// refresh token
	rt := models.RefreshToken{}
	refreshToken := refresh.GenerateRefreshToken(string(user.ID) + deviceAPP.DeviceID + deviceAPP.AppID)
	rt.UserID = user.ID
	rt.DeviceAPPID = deviceAPP.ID
	if tx.Where("device_app_id = ? and user_id = ?", deviceAPP.ID, user.ID).First(&rt).RecordNotFound() {
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
	up := jwt.NewUserPayload(user.ID, deviceAPP.DeviceID, deviceAPP.AppID, m.config.JWT.AccessTokenExpire)
	token, err := jwt.GenerateRSAJWT(up)

	if err != nil {
		return nil, err.Wrapper(errors.New("JWT token generate failed"))
	}

	res["jwt"] = token.String()
	res["refresh_token"] = rt.RefreshToken

	// add operation history
	if err := historyOperation(tx, OperationWechatLogin, user.ID); err != nil {
		return nil, err
	}

	tx.Commit()
	return res, nil
}

func getWechatAccessToken(code string, cfg *config.WechatLoginConfig) (accessToken string, openID string, pe *perror.PlutoError) {
	defer func() {
		var err error
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
		pe = perror.ServerError.Wrapper(err)
	}()
	// get access token
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/oauth2/access_token?appid=%s&secret=%s&code=%s&grant_type=authorization_code",
		*cfg.AppID, *cfg.Secret, code)

	resp, err := http.Get(url)
	if err != nil {
		return "", "", perror.ServerError.Wrapper(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", "", perror.ServerError.Wrapper(err)
	}

	body := make(map[string]interface{})
	if err := json.Unmarshal(b, &body); err != nil {
		return "", "", perror.ServerError.Wrapper(err)
	}

	if resp.StatusCode == http.StatusOK {
		if !strings.Contains(body["scope"].(string), "snsapi_userinfo") {
			return "", "", perror.ServerError.Wrapper(errors.New("Not contain a userinfo scope"))
		}
		return body["access_token"].(string), body["openid"].(string), nil
	}

	if errcode, ok := body["errcode"]; ok {
		// invalid code
		if int(errcode.(float64)) == 40029 {
			return "", "", perror.InvalidWechatCode
		}
		return "", "", perror.ServerError.Wrapper(errors.New(body["errmsg"].(string)))
	}

	return "", "", perror.ServerError.Wrapper(errors.New("Unknow server error"))
}

type wechatUserInfo struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	Sex        string `json:"sex"`
	Province   string `json:"province"`
	City       string `json:"city"`
	Country    string `json:"country"`
	HeadimgURL string `json:"headimgurl"`
	Unionid    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMSG     string `json:"errmsg"`
}

func getWebchatUserInfo(accessToken string, openID string) (info *wechatUserInfo, pe *perror.PlutoError) {

	defer func() {
		var err error
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
		pe = perror.ServerError.Wrapper(err)
	}()
	// get access token
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/userinfo?access_token=%s&openid=%s",
		accessToken, openID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	winfo := wechatUserInfo{}

	if err := json.Unmarshal(b, &winfo); err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	if resp.StatusCode == http.StatusOK {
		return &winfo, nil
	}

	if winfo.ErrMSG != "" {
		// invalid code
		if winfo.ErrCode == 40003 {
			return nil, perror.InvalidWechatCode
		}
		return nil, perror.ServerError.Wrapper(errors.New(winfo.ErrMSG))
	}

	return nil, perror.ServerError.Wrapper(errors.New("Unknow server error"))
}
