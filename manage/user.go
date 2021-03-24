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
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/BillSJC/appleLogin"
	"github.com/volatiletech/sqlboiler/boil"

	"github.com/volatiletech/sqlboiler/queries/qm"

	"github.com/MuShare/pluto/config"
	"github.com/MuShare/pluto/modelexts"

	perror "github.com/MuShare/pluto/datatype/pluto_error"
	"github.com/MuShare/pluto/datatype/request"
	"github.com/MuShare/pluto/models"
	"github.com/MuShare/pluto/utils/avatar"
	"github.com/MuShare/pluto/utils/jwt"
	saltUtil "github.com/MuShare/pluto/utils/salt"
	gjwt "github.com/dgrijalva/jwt-go"
	"google.golang.org/api/oauth2/v2"
)

const (
	MAILLOGIN   = "mail"
	GOOGLELOGIN = "google"
	WECHATLOGIN = "wechat"
	APPLELOGIN  = "apple"
)

func (m *Manager) randomUserName(exec boil.Executor, prefix string) (string, *perror.PlutoError) {
	randomTokenLen := 5
	name := prefix

	for {
		exists, err := models.Users(qm.Where("name = ?", name)).Exists(exec)
		if err != nil {
			return "", perror.ServerError.Wrapper(err)
		}
		if !exists {
			break
		}
		name = fmt.Sprintf("%s_%s", prefix, saltUtil.RandomToken(randomTokenLen))
	}

	return name, nil
}

func (m *Manager) RandomUserName(prefix string) (string, *perror.PlutoError) {
	randomTokenLen := 5
	name := fmt.Sprintf("%s_%s", prefix, saltUtil.RandomToken(randomTokenLen))

	for {
		exists, err := models.Users(qm.Where("name = ?", name)).Exists(m.db)
		if err != nil {
			return "", perror.ServerError.Wrapper(err)
		}
		if !exists {
			break
		}
		name = fmt.Sprintf("%s_%s", prefix, saltUtil.RandomToken(randomTokenLen))
	}

	return name, nil
}

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

	if user.Verified.Bool == false || mailBinding.Verified.Bool == false {
		return nil, perror.MailIsNotVerified
	}

	salt, err := models.Salts(qm.Where("user_id = ?", user.ID)).One(tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, perror.PasswordNotSet
		}
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

	user, err := models.Users(qm.Where("user_id = ?", login.Account)).One(tx)

	if err != nil && err == sql.ErrNoRows {
		return nil, perror.UserIdNotExist
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	if user.Verified.Bool == false {
		return nil, perror.MailIsNotVerified
	}

	salt, err := models.Salts(qm.Where("user_id = ?", user.ID)).One(tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, perror.PasswordNotSet
		}
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

//check userID unique before this method!!
func (m *Manager) newUser(exec boil.Executor, name, avatar, password string, userID *string, verified bool) (*models.User, *perror.PlutoError) {
	user := &models.User{}
	user.Avatar.SetValid(avatar)
	user.Password.SetValid(password)
	user.Name = name
	user.Verified.SetValid(verified)
	if userID != nil {
		userIDExists, err := models.Users(qm.Where("user_id = ?", *userID)).Exists(exec)
		if err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
		if userIDExists {
			return nil, perror.UserIdExists
		}
		user.UserID = *userID
		user.UserIDUpdated = true
	} else {
		newUuid, uuidErr := uuid.NewRandom()
		if uuidErr != nil {
			return nil, perror.ServerError.Wrapper(uuidErr)
		}
		if uuidString := newUuid.String(); uuidString != "" {
			user.UserID = uuidString
			user.UserIDUpdated = false
		} else {
			return nil, perror.ServerError.Wrapper(errors.New("invalid uuid"))
		}
	}
	if err := user.Insert(exec, boil.Infer()); err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	return user, nil
}

func (m *Manager) newBinding(exec boil.Executor, userID uint, mail, loginType, identifyToken string, verified bool) (*models.Binding, *perror.PlutoError) {
	binding := &models.Binding{}
	binding.UserID = userID
	binding.LoginType = loginType
	binding.IdentifyToken = identifyToken
	binding.Mail = mail
	binding.Verified.SetValid(verified)

	if err := binding.Insert(exec, boil.Infer()); err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}
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

	googleBinding, err := models.Bindings(qm.Where("login_type = ? and identify_token = ?", GOOGLELOGIN, info.Sub)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	salt := saltUtil.RandomSalt(info.Sub)

	randomPassword := saltUtil.RandomToken(10)
	encodedPassword, perr := saltUtil.EncodePassword(randomPassword, salt)
	if perr != nil {
		return nil, perr
	}

	namePrefix := ""

	if info.Name == "" {
		namePrefix = "google_user"
	} else {
		namePrefix = info.Name
	}

	name, perr := m.randomUserName(tx, namePrefix)

	if perr != nil {
		return nil, perr
	}

	var user *models.User
	if googleBinding == nil {
		user, perr = m.newUser(tx, name, info.Picture, encodedPassword, nil, true)
		if perr != nil {
			return nil, perr
		}
		googleBinding, perr = m.newBinding(tx, user.ID, info.Email, GOOGLELOGIN, info.Sub, true)
		if perr != nil {
			return nil, perr
		}
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

	identifyToken := info.Unionid
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

	namePrefix := ""

	if info.Nickname == "" {
		namePrefix = "wechat_user"
	} else {
		namePrefix = info.Nickname
	}

	name, perr := m.randomUserName(tx, namePrefix)

	if perr != nil {
		return nil, perr
	}

	var user *models.User
	if wechatBinding == nil {
		user, perr = m.newUser(tx, name, info.HeadimgURL, encodedPassword, nil, true)
		if perr != nil {
			return nil, perr
		}
		wechatBinding, perr = m.newBinding(tx, user.ID, info.Nickname, WECHATLOGIN, info.Unionid, true)
		if perr != nil {
			return nil, perr
		}
	} else {
		wechatBinding.Mail = info.Nickname
		if _, err := wechatBinding.Update(tx, boil.Whitelist("mail")); err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
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
		if err != nil {
			pe = perror.ServerError.Wrapper(err)
		}
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
	Sex        int    `json:"sex"`
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
		if err != nil {
			pe = perror.ServerError.Wrapper(err)
		}
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

	identifyToken := info.Sub
	appleBinding, err := models.Bindings(qm.Where("login_type = ? and identify_token = ?", APPLELOGIN, info.Sub)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	salt := saltUtil.RandomSalt(identifyToken)

	randomPassword := saltUtil.RandomToken(10)
	encodedPassword, perr := saltUtil.EncodePassword(randomPassword, salt)
	if perr != nil {
		return nil, perr
	}

	name, perr := m.randomUserName(tx, "apple_user")

	if perr != nil {
		return nil, perr
	}

	var user *models.User
	if appleBinding == nil {
		user, perr = m.newUser(tx, name, avatarURL, encodedPassword, nil, true)
		if perr != nil {
			return nil, perr
		}
		appleBinding, perr = m.newBinding(tx, user.ID, info.Email, APPLELOGIN, info.Sub, true)
		if perr != nil {
			return nil, perr
		}
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
	_, err := models.Bindings(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).One(m.db)
	if err != nil && err == sql.ErrNoRows {
		return perror.MailNotExist
	} else if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	return nil
}

func (m *Manager) ResetPasswordPage(token string) *perror.PlutoError {

	jwtToken, perr := jwt.VerifyB64RS256JWT(token)
	// token verify failed
	if perr != nil {
		return perr
	}

	prp := &jwt.PasswordResetPayload{}

	if perr := jwtToken.UnmarshalPayload(prp); perr != nil {
		return perr
	}

	if prp.Type != jwt.PASSWORDRESET {
		return perror.InvalidJWTToken
	}

	if time.Now().Unix() > prp.Expire {
		return perror.JWTTokenExpired
	}

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(prp.Mail))
	binding, err := models.Bindings(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).One(m.db)
	if err != nil && err == sql.ErrNoRows {
		return perror.ServerError.Wrapper(errors.New("mail not found"))
	} else if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	user, err := models.Users(qm.Where("id = ?", binding.UserID)).One(m.db)

	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	// user is updated after password reset token is created
	if user.UpdatedAt.Valid && user.UpdatedAt.Time.Unix() > prp.Create {
		return perror.InvalidJWTToken
	}

	return nil
}

func (m *Manager) ResetPassword(token string, rp request.ResetPasswordWeb) *perror.PlutoError {

	jwtToken, perr := jwt.VerifyB64RS256JWT(token)
	if perr != nil {
		return perr
	}

	prp := &jwt.PasswordResetPayload{}

	if perr := jwtToken.UnmarshalPayload(prp); perr != nil {
		return perr
	}

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
	binding, err := models.Bindings(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).One(tx)
	if err != nil && err == sql.ErrNoRows {
		return perror.ServerError.Wrapper(errors.New("mail not found"))
	} else if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	user, err := models.Users(qm.Where("id = ?", binding.UserID)).One(tx)

	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	// user is updated after password reset token is created
	if user.UpdatedAt.Valid && user.UpdatedAt.Time.Unix() > prp.Create {
		return perror.InvalidJWTToken
	}

	salt, err := models.Salts(qm.Where("user_id = ?", user.ID)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return perror.ServerError.Wrapper(err)
	}

	saltString := saltUtil.RandomSalt(prp.Mail)
	if salt == nil {
		salt := models.Salt{}
		salt.Salt = saltString
		salt.UserID = user.ID
		if err := salt.Insert(tx, boil.Infer()); err != nil {
			return perror.ServerError.Wrapper(err)
		}
	} else {
		salt.Salt = saltString
		if _, err := salt.Update(tx, boil.Whitelist("salt")); err != nil {
			return perror.ServerError.Wrapper(err)
		}
	}

	encodedPassword, perr := saltUtil.EncodePassword(rp.Password, saltString)
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

	bindings, err := models.Bindings(qm.Where("user_id = ? and verified = ?", userID, true)).All(m.db)

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	isSaltExsits, err := models.Salts(qm.Where("user_id = ?", user.ID)).Exists(m.db)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	userExt := &modelexts.User{
		User:        user,
		Bindings:    bindings,
		AppID:       accessPayload.AppID,
		PasswordSet: isSaltExsits,
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
	} else if uui.Avatar != "" {
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
	}

	if uui.Name != "" {
		user.Name = uui.Name
	}

	if uui.UserID != "" {
		exists, err := models.Users(qm.Where("user_id = ? and id != ?", uui.UserID, user.ID)).Exists(tx)
		if err != nil {
			return perror.ServerError.Wrapper(err)
		}
		if exists {
			return perror.UserIdExists
		}
		user.UserID = uui.UserID
		user.UserIDUpdated = true
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

func (m *Manager) RegisterWithEmail(register request.MailRegister, admin bool) (*models.User, *perror.PlutoError) {

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
		user := &models.User{
			ID: mailBinding.ID,
		}
		return user, perror.MailIsAlreadyRegister
	}

	var userID *string = nil
	if register.UserID != "" {
		userIDExists, err := models.Users(qm.Where("user_id = ?", register.UserID)).Exists(tx)
		if err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
		if userIDExists {
			return nil, perror.UserIdExists
		}
		userID = &register.UserID
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

	verified := false
	if m.config.Server.SkipRegisterVerifyMail || admin {
		verified = true
	}

	user, perr := m.newUser(tx, register.Name, avatarURL, encodedPassword, userID, verified)
	if perr != nil {
		return nil, perr
	}

	_, perr = m.newBinding(tx, user.ID, register.Mail, MAILLOGIN, identifyToken, verified)

	if perr != nil {
		return nil, perr
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

func (m *Manager) RegisterVerifyMail(rvm request.RegisterVerifyMail) (*models.Binding, *perror.PlutoError) {

	var userMail string
	var binding *models.Binding
	var queryErr error
	if rvm.Mail != "" {
		userMail = rvm.Mail
		identifyToken := b64.RawStdEncoding.EncodeToString([]byte(userMail))
		binding, queryErr = models.Bindings(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).One(m.db)
		if queryErr != nil && queryErr == sql.ErrNoRows {
			return nil, perror.MailNotExist
		} else if queryErr != nil {
			return nil, perror.ServerError.Wrapper(queryErr)
		}
	} else {
		user, userErr := models.Users(qm.Where("user_id = ?", rvm.UserID)).One(m.db)
		if userErr != nil && userErr == sql.ErrNoRows {
			return nil, perror.UserIdNotExist
		} else if userErr != nil {
			return nil, perror.ServerError.Wrapper(userErr)
		}
		binding, queryErr = models.Bindings(qm.Where("login_type = ? and user_id = ?", MAILLOGIN, user.ID)).One(m.db)
		if queryErr != nil && queryErr == sql.ErrNoRows {
			return nil, perror.MailNotExist
		} else if queryErr != nil {
			return nil, perror.ServerError.Wrapper(queryErr)
		}
	}

	if binding.Verified.Bool == true {
		return nil, perror.MailAlreadyVerified
	}

	return binding, nil
}

func (m *Manager) RegisterVerify(token string) *perror.PlutoError {

	jwtToken, perr := jwt.VerifyB64RS256JWT(token)
	if perr != nil {
		return perr
	}

	verifyPayload := &jwt.RegisterVerifyPayload{}
	if perr := jwtToken.UnmarshalPayload(verifyPayload); perr != nil {
		return perr
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

	user.Verified.SetValid(true)

	if _, err := user.Update(tx, boil.Whitelist("verified")); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	binding, err := models.Bindings(qm.Where("user_id = ? and login_type = ?", verifyPayload.UserID, MAILLOGIN)).One(tx)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	if binding.Verified.Bool == true {
		return perror.MailAlreadyVerified
	}

	binding.Verified.SetValid(true)

	if _, err := binding.Update(tx, boil.Whitelist("verified")); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	tx.Commit()

	return nil
}

func (m *Manager) BindMail(binding *request.Binding, accessPayload *jwt.AccessPayload) *perror.PlutoError {
	tx, err := m.db.Begin()

	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	existBinding, err := models.Bindings(qm.Where("user_id = ? and login_type = ?", accessPayload.UserID, MAILLOGIN)).One(tx)

	if err != nil && err != sql.ErrNoRows {
		return perror.ServerError.Wrapper(err)
	}
	// If the existing binding is found and the mail is verified,
	// this mail is not allow to bind again
	if existBinding != nil {
		if existBinding.Verified.Bool == true {
			return perror.BindAlreadyExists
		} else {
			if _, err := existBinding.Delete(tx); err != nil {
				return perror.ServerError.Wrapper(err)
			}
		}
	}

	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(binding.Mail))

	exists, err := models.Bindings(qm.Where("login_type = ? and identify_token = ?", MAILLOGIN, identifyToken)).Exists(tx)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	if exists {
		return perror.BindAlreadyExists
	}

	_, perr := m.newBinding(tx, accessPayload.UserID, binding.Mail, MAILLOGIN, identifyToken, false)
	if perr != nil {
		return perr
	}

	tx.Commit()

	return nil
}

func (m *Manager) BindGoogle(binding *request.Binding, accessPayload *jwt.AccessPayload) *perror.PlutoError {

	info, perr := verifyGoogleIdToken(binding.IDToken)
	if perr != nil {
		return perr
	}

	tx, err := m.db.Begin()

	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	exists, err := models.Bindings(qm.Where("user_id = ? and login_type = ?", accessPayload.UserID, GOOGLELOGIN)).Exists(tx)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	if exists {
		return perror.BindAlreadyExists
	}

	identifyToken := info.Sub

	exists, err = models.Bindings(qm.Where("login_type = ? and identify_token = ?", GOOGLELOGIN, identifyToken)).Exists(tx)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	if exists {
		return perror.BindAlreadyExists
	}

	_, perr = m.newBinding(tx, accessPayload.UserID, info.Email, GOOGLELOGIN, identifyToken, true)
	if perr != nil {
		return perr
	}

	tx.Commit()

	return nil
}

func (m *Manager) BindApple(binding *request.Binding, accessPayload *jwt.AccessPayload) *perror.PlutoError {
	info, perr := getAppleToken(m.config, binding.Code)
	if perr != nil {
		return perr
	}

	tx, err := m.db.Begin()

	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	exists, err := models.Bindings(qm.Where("user_id = ? and login_type = ?", accessPayload.UserID, APPLELOGIN)).Exists(tx)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	if exists {
		return perror.BindAlreadyExists
	}

	identifyToken := info.Sub

	exists, err = models.Bindings(qm.Where("login_type = ? and identify_token = ?", APPLELOGIN, identifyToken)).Exists(tx)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	if exists {
		return perror.BindAlreadyExists
	}

	_, perr = m.newBinding(tx, accessPayload.UserID, info.Email, APPLELOGIN, identifyToken, true)
	if perr != nil {
		return perr
	}

	tx.Commit()

	return nil
}

func (m *Manager) BindWechat(binding *request.Binding, accessPayload *jwt.AccessPayload) *perror.PlutoError {

	accessToken, openID, perr := getWechatAccessToken(binding.Code, m.config.WechatLogin)

	if perr != nil {
		return perr
	}

	info, perr := getWechatUserInfo(accessToken, openID)
	if perr != nil {
		return perr
	}

	tx, err := m.db.Begin()

	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	exists, err := models.Bindings(qm.Where("user_id = ? and login_type = ?", accessPayload.UserID, WECHATLOGIN)).Exists(tx)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	if exists {
		return perror.BindAlreadyExists
	}

	identifyToken := info.Unionid

	exists, err = models.Bindings(qm.Where("login_type = ? and identify_token = ?", WECHATLOGIN, identifyToken)).Exists(tx)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	if exists {
		return perror.BindAlreadyExists
	}

	_, perr = m.newBinding(tx, accessPayload.UserID, info.Nickname, WECHATLOGIN, identifyToken, true)
	if perr != nil {
		return perr
	}

	tx.Commit()

	return nil
}

func (m *Manager) Unbind(ub *request.UnBinding, accessPayload *jwt.AccessPayload) *perror.PlutoError {
	tx, err := m.db.Begin()

	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	bindingCount, err := models.Bindings(qm.Where("user_id = ?", accessPayload.UserID)).Count(tx)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}
	if bindingCount == 1 {
		return perror.UnbindNowAllow
	}

	binding, err := models.Bindings(qm.Where("user_id = ? and login_type = ?", accessPayload.UserID, ub.Type)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		return perror.BindNotExist
	}

	fmt.Println(binding)

	if _, err := binding.Delete(tx); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	tx.Commit()

	return nil
}

func (m *Manager) PublicUserInfo(userID string) (map[string]interface{}, *perror.PlutoError) {

	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil, perror.Forbidden
	}

	user, err := models.Users(qm.Where("id = ?", id)).One(m.db)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.UserNotExist
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	bindings, err := models.Bindings(qm.Where("user_id = ?", id)).All(m.db)

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	userExt := &modelexts.User{
		User:     user,
		Bindings: bindings,
	}

	return userExt.PublicInfo(), nil
}
