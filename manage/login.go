package manage

import (
	"database/sql"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/BillSJC/appleLogin"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/leeif/pluto/config"

	gjwt "github.com/dgrijalva/jwt-go"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/models"
	"github.com/leeif/pluto/utils/avatar"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/refresh"
	saltUtil "github.com/leeif/pluto/utils/salt"
	"google.golang.org/api/oauth2/v2"
)

const (
	MAILLOGIN   = "mail"
	GOOGLELOGIN = "google"
	WECHATLOGIN = "wechat"
	APPLELOGIN  = "apple"
)

func (m *Manager) addDeviceApp(tx *sql.Tx, deviceID, appID string) (*models.DeviceApp, *perror.PlutoError) {
	// insert deviceID and appID into device table
	deviceApp, err := models.DeviceApps(qm.Where("device_id = ? and app_id = ?", deviceID, appID)).One(tx)

	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		deviceApp = &models.DeviceApp{}
		deviceApp.DeviceID = deviceID
		deviceApp.AppID = appID
		if err := deviceApp.Insert(tx, boil.Infer()); err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}
	return deviceApp, nil
}

func (m *Manager) updateRefreshToken(tx *sql.Tx, userID uint, deviceApp *models.DeviceApp) (string, *perror.PlutoError) {
	refreshToken := refresh.GenerateRefreshToken(string(userID) + string(deviceApp.ID))

	rt, err := models.RefreshTokens(qm.Where("user_id = ? and device_app_id = ?", userID, deviceApp.ID)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return "", perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		rt = &models.RefreshToken{}
		rt.DeviceAppID = deviceApp.ID
		rt.UserID = userID
		rt.RefreshToken = refreshToken
		if err := rt.Insert(tx, boil.Infer()); err != nil {
			return "", perror.ServerError.Wrapper(err)
		}
	} else if err == nil {
		rt.RefreshToken = refreshToken
		if _, err := rt.Update(tx, boil.Infer()); err != nil {
			return "", perror.ServerError.Wrapper(err)
		}
	}
	return refreshToken, nil
}

func (m *Manager) EmailLogin(login request.MailLogin) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)

	tx, err := m.db.Begin()

	defer func() {
		tx.Rollback()
	}()

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(login.Mail))
	user, err := models.Users(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).One(tx)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.MailIsNotExsit
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	if user.Verified.Bool == false {
		return nil, perror.MailIsNotVerified
	}

	salt, err := models.Salts(qm.Where("user_id = ?", user.ID)).One(tx)
	if err != nil {
		return nil, perror.ServerError.Wrapper(errors.New("Salt is not found"))
	}

	encodePassword, perr := saltUtil.EncodePassword(login.Password, salt.Salt)

	if perr != nil {
		return nil, perr.Wrapper(errors.New("Password encoding is failed"))
	}

	if user.Password.String != encodePassword {
		return nil, perror.InvalidPassword
	}

	// insert deviceID and appID into device table
	deviceApp, perr := m.addDeviceApp(tx, login.DeviceID, login.AppID)
	if perr != nil {
		return nil, perr
	}

	// update refresh token
	refreshToken, perr := m.updateRefreshToken(tx, user.ID, deviceApp)

	if perr != nil {
		return nil, perr
	}

	// generate jwt token
	up := jwt.NewUserPayload(user.ID, deviceApp.DeviceID, deviceApp.AppID, MAILLOGIN, m.config.JWT.AccessTokenExpire)
	token, perr := jwt.GenerateRSAJWT(up)

	if perr != nil {
		return nil, perr.Wrapper(errors.New("JWT token generate failed"))
	}

	// add operation history
	if err := historyOperation(tx, OperationMailLogin, user.ID); err != nil {
		return nil, err
	}

	res["jwt"] = token.String()
	res["refresh_token"] = refreshToken

	tx.Commit()

	return res, nil
}

func (m *Manager) GoogleLoginMobile(login request.GoogleMobileLogin) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)

	info, perr := verifyGoogleIdToken(login.IDToken)
	if perr != nil {
		return nil, perr
	}

	tx, err := m.db.Begin()
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	user, err := models.Users(qm.Where("login_type = ? and identify_token = ?", GOOGLELOGIN, info.Sub)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err != nil && err == sql.ErrNoRows {
		user = &models.User{}
		user.IdentifyToken = info.Sub
		user.LoginType = GOOGLELOGIN
		user.Avatar.SetValid(info.Picture)
		user.Name = info.Name
		user.Verified.SetValid(true)
		if err := user.Insert(tx, boil.Infer()); err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	// insert deviceID and appID into device table
	deviceApp, perr := m.addDeviceApp(tx, login.DeviceID, login.AppID)
	if perr != nil {
		return nil, perr
	}

	// update refresh token
	refreshToken, perr := m.updateRefreshToken(tx, user.ID, deviceApp)

	if perr != nil {
		return nil, perr
	}

	// generate jwt token
	up := jwt.NewUserPayload(user.ID, deviceApp.DeviceID, deviceApp.AppID, GOOGLELOGIN, m.config.JWT.AccessTokenExpire)
	token, perr := jwt.GenerateRSAJWT(up)

	if perr != nil {
		return nil, perr.Wrapper(errors.New("JWT token generate failed"))
	}

	// add operation history
	if err := historyOperation(tx, OperationGoogleLogin, user.ID); err != nil {
		return nil, err
	}

	res["jwt"] = token.String()
	res["refresh_token"] = refreshToken

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

func (m *Manager) WechatLoginMobile(login request.WechatMobileLogin) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)

	accessToken, openID, perr := getWechatAccessToken(login.Code, m.config.WechatLogin)

	if perr != nil {
		return nil, perr
	}

	info, perr := getWechatUserInfo(accessToken, openID)
	if perr != nil {
		return nil, perr
	}

	tx, err := m.db.Begin()
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	user, err := models.Users(qm.Where("login_type = ? and identify_token = ?", WECHATLOGIN, info.Unionid)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err != nil && err == sql.ErrNoRows {
		user = &models.User{}
		user.IdentifyToken = info.Unionid
		user.LoginType = WECHATLOGIN
		user.Avatar.SetValid(info.HeadimgURL)
		user.Name = info.Nickname
		user.Verified.SetValid(true)
		if err := user.Insert(tx, boil.Infer()); err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	// insert deviceID and appID into device table
	deviceApp, perr := m.addDeviceApp(tx, login.DeviceID, login.AppID)
	if perr != nil {
		return nil, perr
	}

	// update refresh token
	refreshToken, perr := m.updateRefreshToken(tx, user.ID, deviceApp)

	if perr != nil {
		return nil, perr
	}

	// generate jwt token
	up := jwt.NewUserPayload(user.ID, deviceApp.DeviceID, deviceApp.AppID, WECHATLOGIN, m.config.JWT.AccessTokenExpire)
	token, perr := jwt.GenerateRSAJWT(up)

	if perr != nil {
		return nil, perr.Wrapper(errors.New("JWT token generate failed"))
	}

	// add operation history
	if err := historyOperation(tx, OperationWechatLogin, user.ID); err != nil {
		return nil, err
	}

	res["jwt"] = token.String()
	res["refresh_token"] = refreshToken

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

func getWechatUserInfo(accessToken string, openID string) (info *wechatUserInfo, pe *perror.PlutoError) {

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

func (m *Manager) AppleLoginMobile(login request.AppleMobileLogin) (map[string]string, *perror.PlutoError) {
	res := make(map[string]string)

	info, perr := getAppleToken(m.config, login.Code)

	if perr != nil {
		return nil, perr
	}

	tx, err := m.db.Begin()
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	user, err := models.Users(qm.Where("login_type = ? and identify_token = ?", APPLELOGIN, info.Sub)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err != nil && err == sql.ErrNoRows {
		user = &models.User{}
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
			m.logger.Warn(perr.LogError)
		} else {
			user.Avatar.SetValid(remoteURL)
		}

		user.IdentifyToken = info.Sub
		user.Mail.SetValid(info.Email)
		user.LoginType = APPLELOGIN
		user.Name = login.Name
		user.Verified.SetValid(true)
		if err := user.Insert(tx, boil.Infer()); err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	// insert deviceID and appID into device table
	deviceApp, perr := m.addDeviceApp(tx, login.DeviceID, login.AppID)
	if perr != nil {
		return nil, perr
	}

	// update refresh token
	refreshToken, perr := m.updateRefreshToken(tx, user.ID, deviceApp)

	if perr != nil {
		return nil, perr
	}

	// generate jwt token
	up := jwt.NewUserPayload(user.ID, deviceApp.DeviceID, deviceApp.AppID, APPLELOGIN, m.config.JWT.AccessTokenExpire)
	token, perr := jwt.GenerateRSAJWT(up)

	if perr != nil {
		return nil, perr.Wrapper(errors.New("JWT token generate failed"))
	}

	// add operation history
	if err := historyOperation(tx, OperationAppleLogin, user.ID); err != nil {
		return nil, err
	}

	res["jwt"] = token.String()
	res["refresh_token"] = refreshToken

	tx.Commit()
	return res, nil
}

type appleIdTokenInfo struct {
	Iss           string `json:"iss"`
	Aud           string `json:"aud"`
	Exp           int64  `json:"exp"`
	Iat           int64  `json:"iat"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified string `json:"email_verified"`
	AuthTime      int64  `json:"auth_time"`
	gjwt.StandardClaims
}

func getAppleToken(cfg *config.Config, code string) (*appleIdTokenInfo, *perror.PlutoError) {
	fmt.Println(cfg.AppleLogin.BundleID)
	a := appleLogin.InitAppleConfig(
		cfg.AppleLogin.TeamID,
		cfg.AppleLogin.BundleID,
		cfg.AppleLogin.RedirectURL,
		cfg.AppleLogin.KeyID,
	)

	err := a.LoadP8CertByFile(cfg.AppleLogin.P8CertFile)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	token, err := a.GetAppleToken(code, 120)
	if err != nil {
		return nil, perror.InvalidAppleIDToken.Wrapper(err)
	}

	info, perr := parseAppleIDToken(token.IDToken)
	if perr != nil {
		return nil, perr
	}

	if info.Aud != cfg.AppleLogin.BundleID {
		return nil, perror.InvalidAppleIDToken
	}
	return info, nil
}

func parseAppleIDToken(idToken string) (*appleIdTokenInfo, *perror.PlutoError) {
	parser := gjwt.Parser{}
	token, _, err := parser.ParseUnverified(idToken, &appleIdTokenInfo{})
	if err != nil {
		return nil, perror.InvalidAppleIDToken.Wrapper(err)
	}
	info, ok := token.Claims.(*appleIdTokenInfo)
	if !ok {
		return nil, perror.InvalidAppleIDToken
	}
	return info, nil
}
