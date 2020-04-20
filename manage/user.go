package manage

import (
	"database/sql"
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/BillSJC/appleLogin"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/modelexts"

	gjwt "github.com/dgrijalva/jwt-go"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/models"
	"github.com/leeif/pluto/utils/avatar"
	"github.com/leeif/pluto/utils/jwt"
	saltUtil "github.com/leeif/pluto/utils/salt"
	"google.golang.org/api/oauth2/v2"
)

const (
	MAILLOGIN   = "mail"
	GOOGLELOGIN = "google"
	WECHATLOGIN = "wechat"
	APPLELOGIN  = "apple"
)

func (m *Manager) MailPasswordLogin(login request.PasswordLogin) (*GrantResult, *perror.PlutoError) {
	tx, err := m.db.Begin()

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(login.Account))
	mailBinding, err := models.Bindings(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).One(tx)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.MailNotExist
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	user, err := models.Users(qm.Where("id = ?", mailBinding.UserID)).One(tx)

	if err != nil {
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

	scopes, perr := getUserDefaultScopes(tx, user.ID, login.AppID)
	if perr != nil {
		return nil, perr
	}

	grantResult, perr := m.loginWithAppName(tx, user.ID, login.DeviceID, login.AppID, strings.Join(scopes, ","))
	if perr != nil {
		return nil, perr
	}

	tx.Commit()

	return grantResult, nil
}

func (m *Manager) NamePasswordLogin(login request.PasswordLogin) (*GrantResult, *perror.PlutoError) {
	tx, err := m.db.Begin()

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	user, err := models.Users(qm.Where("name = ?", login.Account)).One(tx)

	if err != nil {
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

	scopes, perr := getUserDefaultScopes(tx, user.ID, login.AppID)
	if perr != nil {
		return nil, perr
	}

	grantResult, perr := m.loginWithAppName(tx, user.ID, login.DeviceID, login.AppID, strings.Join(scopes, ","))
	if perr != nil {
		return nil, perr
	}

	tx.Commit()

	return grantResult, nil
}

func (m *Manager) newUser(exec boil.Executor, name, avatar, password string) (*models.User, *perror.PlutoError) {
	user := &models.User{}
	user.Avatar.SetValid(avatar)
	user.Password.SetValid(password)
	user.Name = name
	user.Verified.SetValid(true)
	if err := user.Insert(exec, boil.Infer()); err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	return user, nil
}

func (m *Manager) newBinding(exec boil.Executor, userID uint, mail, loginType, identifyToken string) (*models.Binding, *perror.PlutoError) {
	binding := &models.Binding{}
	binding.UserID = userID
	binding.LoginType = loginType
	binding.IdentifyToken = identifyToken
	binding.Mail = mail
	return binding, nil
}

func (m *Manager) GoogleLoginMobile(login request.GoogleMobileLogin) (*GrantResult, *perror.PlutoError) {
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

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(info.Sub))
	googleBinding, err := models.Bindings(qm.Where("login_type = ? and identify_token = ?", GOOGLELOGIN, identifyToken)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	salt := saltUtil.RandomSalt(identifyToken)

	randomPassword := saltUtil.RandomToken(10)
	encodedPassword, perr := saltUtil.EncodePassword(randomPassword, salt)
	if perr != nil {
		return nil, perr
	}

	var user *models.User
	if googleBinding != nil {
		user, perr = m.newUser(tx, info.Name, info.Picture, encodedPassword)
		if perr != nil {
			return nil, perr
		}
		googleBinding, perr = m.newBinding(tx, user.ID, info.Email, GOOGLELOGIN, info.Sub)
	} else {
		googleBinding.Mail = info.Email
		if _, err := googleBinding.Update(tx, boil.Infer()); err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
		user, err = models.Users(qm.Where("id = ?", googleBinding.UserID)).One(tx)
		if err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	scopes, perr := getUserDefaultScopes(tx, user.ID, login.AppID)
	if perr != nil {
		return nil, perr
	}

	grantResult, perr := m.loginWithAppName(tx, user.ID, login.DeviceID, login.AppID, strings.Join(scopes, ","))
	if perr != nil {
		return nil, perr
	}

	tx.Commit()

	return grantResult, nil
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

func (m *Manager) WechatLoginMobile(login request.WechatMobileLogin) (*GrantResult, *perror.PlutoError) {
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

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(info.Unionid))
	wechatBinding, err := models.Bindings(qm.Where("login_type = ? and identify_token = ?", WECHATLOGIN, identifyToken)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	salt := saltUtil.RandomSalt(identifyToken)

	randomPassword := saltUtil.RandomToken(10)
	encodedPassword, perr := saltUtil.EncodePassword(randomPassword, salt)
	if perr != nil {
		return nil, perr
	}

	var user *models.User
	if wechatBinding != nil {
		user, perr = m.newUser(tx, info.Nickname, info.HeadimgURL, encodedPassword)
		if perr != nil {
			return nil, perr
		}
		wechatBinding, perr = m.newBinding(tx, user.ID, "", WECHATLOGIN, info.Unionid)
	} else {
		user, err = models.Users(qm.Where("id = ?", wechatBinding.UserID)).One(tx)
		if err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	scopes, perr := getUserDefaultScopes(tx, user.ID, login.AppID)
	if perr != nil {
		return nil, perr
	}

	grantResult, perr := m.loginWithAppName(tx, user.ID, login.DeviceID, login.AppID, strings.Join(scopes, ","))
	if perr != nil {
		return nil, perr
	}

	tx.Commit()
	return grantResult, nil
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

func (m *Manager) AppleLoginMobile(login request.AppleMobileLogin) (*GrantResult, *perror.PlutoError) {
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

	ag := avatar.AvatarGen{}
	avatarReader, perr := ag.GenFromGravatar()
	if perr != nil {
		return nil, perr
	}

	avatarURL := ""
	as := avatar.NewAvatarSaver(m.config)
	remoteURL, perr := as.SaveAvatarImageInOSS(avatarReader)
	if perr != nil {
		avatarURL = avatarReader.OriginURL
		m.logger.Warn(perr.LogError)
	} else {
		avatarURL = remoteURL
	}

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(info.Sub))
	appleBinding, err := models.Bindings(qm.Where("login_type = ? and identify_token = ?", WECHATLOGIN, info.Sub)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	salt := saltUtil.RandomSalt(identifyToken)

	randomPassword := saltUtil.RandomToken(10)
	encodedPassword, perr := saltUtil.EncodePassword(randomPassword, salt)
	if perr != nil {
		return nil, perr
	}

	var user *models.User
	if appleBinding != nil {
		user, perr = m.newUser(tx, login.Name, avatarURL, encodedPassword)
		if perr != nil {
			return nil, perr
		}
		appleBinding, perr = m.newBinding(tx, user.ID, info.Email, APPLELOGIN, info.Sub)
	} else {
		user, err = models.Users(qm.Where("id = ?", appleBinding.UserID)).One(tx)
		if err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	scopes, perr := getUserDefaultScopes(tx, user.ID, login.AppID)
	if perr != nil {
		return nil, perr
	}

	grantResult, perr := m.loginWithAppName(tx, user.ID, login.DeviceID, login.AppID, strings.Join(scopes, ","))
	if perr != nil {
		return nil, perr
	}

	tx.Commit()
	return grantResult, nil
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

func (m *Manager) ResetPasswordMail(rpm request.ResetPasswordMail) *perror.PlutoError {

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(rpm.Mail))
	_, err := models.Users(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).One(m.db)
	if err != nil && err == sql.ErrNoRows {
		return perror.MailNotExist
	} else if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	return nil
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
		return perror.InvalidJWTToken
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
		return perror.InvalidJWTToken
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
		return perror.InvalidJWTToken
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
		return perror.InvalidJWTToken
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

func (m *Manager) UserInfo(userID uint, accessPayload *jwt.AccessPayload) (*modelexts.User, *perror.PlutoError) {

	user, err := models.Users(qm.Where("id = ?", userID)).One(m.db)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(errors.New("user not found id: " + string(accessPayload.UserID)))
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	role, perr := getUserRole(m.db, userID, accessPayload.AppID)
	if perr != nil {
		return nil, perr
	}

	userExt := &modelexts.User{
		User: user,
	}

	if role != nil {
		userExt.Role = role.Name
	}

	return userExt, nil
}

func (m *Manager) UpdateUserInfo(accessPayload *jwt.AccessPayload, uui request.UpdateUserInfo) *perror.PlutoError {

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

func (m *Manager) RegisterWithEmail(register request.MailRegister) (*models.User, *perror.PlutoError) {

	tx, err := m.db.Begin()
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}
	defer func() {
		tx.Rollback()
	}()

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(register.Mail))
	mailBinding, err := models.Bindings(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	if mailBinding != nil {
		return nil, perror.MailIsAlreadyRegister
	}

	salt := saltUtil.RandomSalt(identifyToken)

	encodedPassword, perr := saltUtil.EncodePassword(register.Password, salt)
	if perr != nil {
		return nil, perr
	}

	avatarURL := ""
	// get a random avatar
	ag := avatar.AvatarGen{}
	avatarReader, perr := ag.GenFromGravatar()
	if perr != nil {
		return nil, perr
	}

	as := avatar.NewAvatarSaver(m.config)
	remoteURL, perr := as.SaveAvatarImageInOSS(avatarReader)
	if perr != nil {
		avatarURL = avatarReader.OriginURL
		m.logger.Error(perr.LogError)
	} else {
		avatarURL = remoteURL
	}

	user, perr := m.newUser(tx, register.Name, avatarURL, encodedPassword)
	if perr != nil {
		return nil, perr
	}

	_, perr = m.newBinding(tx, user.ID, register.Mail, MAILLOGIN, identifyToken)

	if perr != nil {
		return nil, perr
	}

	if m.config.Server.SkipRegisterVerifyMail {
		user.Verified.SetValid(true)
		if err := user.Insert(tx, boil.Infer()); err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
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
