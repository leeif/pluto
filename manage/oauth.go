package manage

import (
	"database/sql"
	b64 "encoding/base64"
	"errors"
	"net/url"
	"time"

	"github.com/RichardKnop/uuid"
	"github.com/leeif/pluto/datatype/pluto_error"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/models"
	"github.com/leeif/pluto/utils/general"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/refresh"
	saltUtil "github.com/leeif/pluto/utils/salt"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

const (
	ClientPend     = "pend"
	ClientApproved = "approved"
	ClientDenied   = "denied"
)

type GrantResult struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Type         string `json:"type"`
}

type AuthorizeResult struct {
	Code        string
	State       string
	AccessToken string
	ExpireAt    string
	Type        string
	Scopes      string
}

func (m *Manager) getUser(exec boil.Executor, id uint) (*models.User, *perror.PlutoError) {
	user, err := models.Users(qm.Where("id=?", id)).One(exec)

	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		return nil, perror.UserNotExist
	}

	return user, nil
}

func (m *Manager) getApplication(exec boil.Executor, application string) (*models.Application, *perror.PlutoError) {
	app, err := models.Applications(qm.Where("name=?", application)).One(exec)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		return nil, perror.ApplicationNotExist
	}

	return app, nil
}

func (m *Manager) getApplicationByID(exec boil.Executor, id uint) (*models.Application, *perror.PlutoError) {
	app, err := models.Applications(qm.Where("id=?", id)).One(exec)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		return nil, perror.ApplicationNotExist
	}

	return app, nil
}

func (m *Manager) getClient(exec boil.Executor, id uint) (*models.OauthClient, *perror.PlutoError) {
	client, err := models.OauthClients(qm.Where("id=?", id)).One(exec)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		return nil, perror.OAuthClientIDOrSecretNotFound
	}

	return client, nil
}

func (m *Manager) getClientByKey(exec boil.Executor, key string) (*models.OauthClient, *perror.PlutoError) {
	client, err := models.OauthClients(qm.Where("`key`=?", key)).One(exec)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		return nil, perror.OAuthClientIDOrSecretNotFound
	}

	return client, nil
}

func (m *Manager) GetClientByKey(key string) (*models.OauthClient, *perror.PlutoError) {
	client, err := models.OauthClients(qm.Where("`key`=?", key)).One(m.db)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		return nil, perror.OAuthClientIDOrSecretNotFound
	}

	return client, nil
}

func (m *Manager) authClient(exec boil.Executor, clientID, secret string) (*models.OauthClient, *perror.PlutoError) {
	client, err := models.OauthClients(qm.Where("`key`=?", clientID)).One(exec)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		return nil, perror.OAuthClientIDOrSecretNotFound
	}

	if client.Status != ClientApproved {
		return nil, perror.OAuthInvalidClient
	}

	if err := general.VerifyPassword(client.Secret, secret); err != nil {
		return nil, perror.OAuthInvalidClient
	}

	return client, nil
}

func (m *Manager) AuthorizationCodeGrant(ot *request.OAuthTokens) (*GrantResult, *perror.PlutoError) {
	tx, err := m.db.Begin()

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	client, perr := m.authClient(tx, ot.ClientID, ot.ClientSecret)
	if perr != nil {
		return nil, perr
	}

	authorizationCode, err := models.OauthAuthorizationCodes(qm.Where("code = ? and client_id = ?", ot.Code, client.ID)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		return nil, perror.OAuthAuthorizationCodeNotFound
	}

	if authorizationCode.RedirectURI != ot.RedirectURI {
		return nil, perror.OAuthInvalidRedirectURI
	}

	if time.Now().After(authorizationCode.ExpireAt) {
		return nil, perror.OAuthAuthorizationCodeExpired
	}

	if ot.DeviceID == "" {
		ot.DeviceID = "oauth"
	}

	// login
	grantResult, perr := m.loginWithAppID(tx, authorizationCode.UserID, ot.DeviceID, authorizationCode.AppID, authorizationCode.Scopes)
	if perr != nil {
		return nil, perr
	}

	// delete authorizationCode
	if _, err := authorizationCode.Delete(tx); err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	tx.Commit()

	return grantResult, nil
}

func (m *Manager) PasswordGrant(ot *request.OAuthTokens) (*GrantResult, *perror.PlutoError) {
	tx, err := m.db.Begin()

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	if _, perr := m.authClient(tx, ot.ClientID, ot.ClientSecret); perr != nil {
		return nil, perr
	}

	user, perr := m.authPassword(tx, ot.User, ot.Password)
	if perr != nil {
		return nil, perr
	}

	scopes, perr := getValidScopes(tx, ot.Scopes, user.ID, ot.AppID)
	if perr != nil {
		return nil, perr
	}

	if ot.DeviceID == "" {
		ot.DeviceID = "oauth"
	}

	grantResult, perr := m.loginWithAppName(tx, user.ID, ot.DeviceID, ot.AppID, scopes)
	if perr != nil {
		return nil, perr
	}

	tx.Commit()

	return grantResult, nil
}

func (m *Manager) ClientCredentialGrant(ot *request.OAuthTokens) (*GrantResult, *perror.PlutoError) {

	application, perr := m.getApplication(m.db, ot.AppID)

	if perr != nil {
		return nil, perr
	}

	if _, perr := m.authClient(m.db, ot.ClientID, ot.ClientSecret); perr != nil {
		return nil, perr
	}

	// empty user, scopes, refreshToken
	return m.grantToken(0, "", "", application.Name, m.config.Token.AccessTokenExpire)
}

func (m *Manager) RefreshTokenGrant(ot *request.OAuthTokens) (*GrantResult, *perror.PlutoError) {

	tx, err := m.db.Begin()

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	if _, perr := m.authClient(tx, ot.ClientID, ot.ClientSecret); perr != nil {
		return nil, perr
	}

	rt, perr := m.authRefreshToken(tx, ot.RefreshToken)
	if perr != nil {
		return nil, perr
	}

	deviceApp, err := models.DeviceApps(qm.Where("id = ?", rt.DeviceAppID)).One(tx)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	application, err := models.Applications(qm.Where("id = ?", deviceApp.AppID)).One(tx)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	if application.Name != ot.AppID {
		return nil, perror.InvalidApplication
	}

	scopes, perr := getValidScopes(tx, ot.Scopes, rt.UserID, application.Name)

	if perr != nil {
		return nil, perr
	}

	if scopes == "" && !rt.Scopes.IsZero() {
		scopes = rt.Scopes.String
	}

	if perr := m.updateRefreshToken(tx, rt, scopes); perr != nil {
		return nil, perr
	}

	grantResult, perr := m.grantToken(rt.UserID, rt.RefreshToken, scopes, ot.AppID, m.config.Token.AccessTokenExpire)
	if perr != nil {
		return nil, perr
	}

	tx.Commit()

	return grantResult, nil
}

func (m *Manager) authRefreshToken(exec boil.Executor, refreshToken string) (*models.RefreshToken, *perror.PlutoError) {

	rt, err := models.RefreshTokens(qm.Where("refresh_token = ?", refreshToken)).One(exec)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.InvalidRefreshToken
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	return rt, nil
}

func (m *Manager) authPassword(exec boil.Executor, mail, password string) (*models.User, *perror.PlutoError) {
	identifyToken := b64.RawStdEncoding.EncodeToString([]byte(mail))
	user, err := models.Users(qm.Where("identify_token = ?", identifyToken)).One(exec)
	if err != nil && err == sql.ErrNoRows {
		return nil, perror.MailNotExist
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	if user.Verified.Bool == false {
		return nil, perror.MailIsNotVerified
	}

	salt, err := models.Salts(qm.Where("user_id = ?", user.ID)).One(exec)
	if err != nil {
		return nil, perror.ServerError.Wrapper(errors.New("Salt is not found"))
	}

	encodePassword, perr := saltUtil.EncodePassword(password, salt.Salt)

	if perr != nil {
		return nil, perr.Wrapper(errors.New("Password encoding is failed"))
	}

	if user.Password.String != encodePassword {
		return nil, perror.InvalidPassword
	}

	return user, nil
}

func (m *Manager) loginWithAppName(exec boil.Executor, userID uint, deviceID, appName string, scopes string) (*GrantResult, *perror.PlutoError) {

	application, perr := m.getApplication(exec, appName)
	if perr != nil {
		return nil, perr
	}

	// insert deviceID and appID into device table
	deviceApp, perr := m.getDeviceApp(exec, deviceID, application)
	if perr != nil {
		return nil, perr
	}

	// create refresh token
	refreshToken, perr := m.newRefreshToken(exec, userID, deviceApp, scopes)

	if perr != nil {
		return nil, perr
	}

	// grant access token
	return m.grantToken(userID, refreshToken.RefreshToken, scopes, application.Name, m.config.Token.AccessTokenExpire)
}

func (m *Manager) loginWithAppID(exec boil.Executor, userID uint, deviceID string, appID uint, scopes string) (*GrantResult, *perror.PlutoError) {

	application, perr := m.getApplicationByID(exec, appID)
	if perr != nil {
		return nil, perr
	}

	// insert deviceID and appID into device table
	deviceApp, perr := m.getDeviceApp(exec, deviceID, application)
	if perr != nil {
		return nil, perr
	}

	// update refresh token
	refreshToken, perr := m.newRefreshToken(exec, userID, deviceApp, scopes)

	if perr != nil {
		return nil, perr
	}

	// grant access token
	return m.grantToken(userID, refreshToken.RefreshToken, scopes, application.Name, m.config.Token.AccessTokenExpire)
}

func (m *Manager) newRefreshToken(exec boil.Executor, userID uint, deviceApp *models.DeviceApp, scopes string) (*models.RefreshToken, *perror.PlutoError) {
	refreshToken := refresh.GenerateRefreshToken(string(userID) + string(deviceApp.ID))

	rt := &models.RefreshToken{}
	rt.DeviceAppID = deviceApp.ID
	rt.UserID = userID
	rt.RefreshToken = refreshToken
	rt.Scopes.SetValid(scopes)
	if err := rt.Insert(exec, boil.Infer()); err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	return rt, nil
}

func (m *Manager) updateRefreshToken(exec boil.Executor, rt *models.RefreshToken, scopes string) *perror.PlutoError {
	newToken := refresh.GenerateRefreshToken(string(rt.UserID) + string(rt.DeviceAppID))
	rt.RefreshToken = newToken
	rt.Scopes.SetValid(scopes)

	if _, err := rt.Update(exec, boil.Infer()); err != nil {
		return perror.ServerError.Wrapper(err)
	}
	return nil
}

func (m *Manager) getDeviceApp(exec boil.Executor, deviceID string, application *models.Application) (*models.DeviceApp, *perror.PlutoError) {

	// insert deviceID and appID into device table
	deviceApp, err := models.DeviceApps(qm.Where("device_id = ? and app_id = ?", deviceID, application.ID)).One(exec)

	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		deviceApp = &models.DeviceApp{}
		deviceApp.DeviceID = deviceID
		deviceApp.AppID = application.ID
		if err := deviceApp.Insert(exec, boil.Infer()); err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}
	return deviceApp, nil
}

func (m *Manager) grantToken(userID uint, refreshToken, scopes, appID string, expire int64) (*GrantResult, *perror.PlutoError) {
	// generate jwt token
	up := jwt.NewAccessPayload(userID, scopes, appID, expire)
	accessToken, perr := jwt.GenerateRSAJWT(up)

	if perr != nil {
		return nil, perr.Wrapper(errors.New("JWT token generate failed"))
	}

	grantResult := &GrantResult{
		Type:         "Bearer",
		AccessToken:  accessToken.String(),
		RefreshToken: refreshToken,
	}

	return grantResult, nil
}

func (m *Manager) GrantAuthorizationCode(oa *request.OAuthAuthorize, accessPayload *jwt.AccessPayload) (*AuthorizeResult, *url.URL, *pluto_error.PlutoError) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	client, perr := m.getClientByKey(tx, oa.ClientID)
	if perr != nil {
		return nil, nil, perr
	}

	if client.Status != ClientApproved {
		return nil, nil, perror.OAuthInvalidClient
	}

	if oa.RedirectURI == "" {
		oa.RedirectURI = client.RedirectURI
	}

	redirectURI, err := url.Parse(oa.RedirectURI)

	if err != nil {
		return nil, nil, perror.ServerError.Wrapper(err)
	}

	application, perr := m.getApplication(tx, oa.AppID)

	if perr != nil {
		return nil, redirectURI, perr
	}

	user, perr := m.getUser(tx, accessPayload.UserID)

	scopes, perr := getValidScopes(tx, oa.Scopes, user.ID, oa.AppID)

	if perr != nil {
		return nil, redirectURI, perr
	}

	code := &models.OauthAuthorizationCode{}

	code.ClientID = client.ID
	code.AppID = application.ID
	code.UserID = user.ID
	code.Code = uuid.New()
	code.RedirectURI = redirectURI.String()
	code.ExpireAt = time.Now().Add(time.Duration(m.config.OAuth.AuthorizeCodeExpire) * time.Second)
	code.Scopes = scopes

	if err := code.Insert(tx, boil.Infer()); err != nil {
		return nil, redirectURI, perror.ServerError.Wrapper(err)
	}

	tx.Commit()

	return &AuthorizeResult{
		Code:  code.Code,
		State: oa.State,
	}, redirectURI, nil
}

func (m *Manager) GrantAccessToken(oa *request.OAuthAuthorize, accessPayload *jwt.AccessPayload) (*AuthorizeResult, *url.URL, *pluto_error.PlutoError) {

	client, perr := m.getClientByKey(m.db, oa.ClientID)
	if perr != nil {
		return nil, nil, perr
	}

	if client.Status != ClientApproved {
		return nil, nil, perror.OAuthInvalidClient
	}

	if oa.RedirectURI == "" {
		oa.RedirectURI = client.RedirectURI
	}

	redirectURI, err := url.Parse(oa.RedirectURI)

	if err != nil {
		return nil, nil, perror.ServerError.Wrapper(err)
	}

	if accessPayload.AppID != oa.AppID {
		return nil, redirectURI, perror.InvalidApplication
	}

	application, perr := m.getApplication(m.db, oa.AppID)

	if perr != nil {
		return nil, redirectURI, perr
	}

	user, perr := m.getUser(m.db, accessPayload.UserID)
	if perr != nil {
		return nil, redirectURI, perr
	}

	scopes, perr := getValidScopes(m.db, oa.Scopes, user.ID, oa.AppID)

	if perr != nil {
		return nil, redirectURI, perr
	}

	if oa.LifeTime == 0 {
		oa.LifeTime = m.config.Token.AccessTokenExpire
	}

	grantResult, perr := m.grantToken(user.ID, "", scopes, application.Name, oa.LifeTime)

	if perr != nil {
		return nil, redirectURI, perr
	}

	return &AuthorizeResult{
		AccessToken: grantResult.AccessToken,
		Type:        grantResult.Type,
		State:       oa.State,
		Scopes:      scopes,
	}, redirectURI, nil
}

func (m *Manager) OAuthCreateClient(occ *request.OAuthCreateClient) (*models.OauthClient, *perror.PlutoError) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	if _, err := models.OauthClients(qm.Where("`key`=?", occ.Key)).One(tx); err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == nil {
		return nil, perror.OAuthClientExist
	}

	client := &models.OauthClient{}
	client.Key = occ.Key
	secretHash, err := general.HashPassword(occ.Secret)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}
	client.Secret = string(secretHash)
	client.Status = ClientPend
	client.RedirectURI = occ.RedirectURI

	if err := client.Insert(tx, boil.Infer()); err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	tx.Commit()

	return client, nil
}

func (m *Manager) OAuthApproveClient(occ *request.OAuthApproveClient) (*models.OauthClient, *perror.PlutoError) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	client, err := models.OauthClients(qm.Where("`key`=?", occ.Key)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		return nil, perror.OAuthClientIDOrSecretNotFound
	}

	client.Status = ClientApproved

	if _, err := client.Update(tx, boil.Infer()); err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	return client, nil
}

func (m *Manager) OAuthDenyClient(occ *request.OAuthApproveClient) (*models.OauthClient, *perror.PlutoError) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	client, err := models.OauthClients(qm.Where("`key`=?", occ.Key)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		return nil, perror.OAuthClientIDOrSecretNotFound
	}

	client.Status = ClientDenied

	if _, err := client.Update(tx, boil.Infer()); err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	return client, nil
}
