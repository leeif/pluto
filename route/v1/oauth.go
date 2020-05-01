package v1

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"time"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/manage"
	"github.com/leeif/pluto/utils/general"
	routeUtils "github.com/leeif/pluto/utils/route"
)

func (router *Router) fetchClient(r *http.Request) (string, string, *perror.PlutoError) {
	clientID, secret, ok := r.BasicAuth()
	if !ok {
		return "", "", perror.OAuthClientIDOrSecretNotFound
	}

	return clientID, secret, nil
}

func (router *Router) OAuthTokens(w http.ResponseWriter, r *http.Request) *perror.PlutoError {

	tokens := &request.OAuthTokens{}

	if err := routeUtils.GetRequestData(r, tokens); err != nil {
		return err
	}

	type grantHandler func(*request.OAuthTokens) (*manage.GrantResult, *perror.PlutoError)

	var handler grantHandler
	switch tokens.GrantType {
	case "authorization_code":
		if !tokens.ValidateAuthorizationCode() {
			return perror.BadRequest
		}
		handler = router.manager.AuthorizationCodeGrant
	case "password":
		if !tokens.ValidatePasswordGrant() {
			return perror.BadRequest
		}
		handler = router.manager.PasswordGrant
	case "client_credentials":
		if !tokens.ValidateClientCredentials() {
			return perror.BadRequest
		}
		handler = router.manager.ClientCredentialGrant
	case "refresh_token":
		if !tokens.ValidateRefreshToken() {
			return perror.BadRequest
		}
		handler = router.manager.RefreshTokenGrant
	default:
		return perror.OAuthInvalidGrantType
	}

	clientID, secret, perr := router.fetchClient(r)
	if perr != nil {
		return perr
	}

	tokens.ClientID = clientID
	tokens.ClientSecret = secret

	res, perr := handler(tokens)
	if perr != nil {
		return perr
	}

	if err := routeUtils.ResponseOK(res, w); err != nil {
		return err
	}

	return nil
}

func (router *Router) OAuthAuthorize(w http.ResponseWriter, r *http.Request) *perror.PlutoError {

	accessPayload, perr := routeUtils.GetAccessPayload(r)

	if perr != nil {
		return perr
	}

	authorize := &request.OAuthAuthorize{}

	if perr := routeUtils.GetRequestData(r, authorize); perr != nil {
		return perr
	}

	if authorize.ResponseType == "code" {
		// Create a new authorization code
		authorizeResult, redirectURI, perr := router.manager.GrantAuthorizationCode(authorize, accessPayload)
		if perr != nil && redirectURI != nil {
			router.logger.Debug(perr)
			routeUtils.ErrorRedirect(w, r, redirectURI, perr.PlutoCode, authorize.State, authorize.ResponseType)
			return nil
		} else if perr != nil {
			return perr
		}

		query := redirectURI.Query()
		query.Set("code", authorizeResult.Code)
		if authorize.State != "" {
			query.Set("state", authorizeResult.State)
		}

		routeUtils.RedirectWithQueryString(redirectURI.String(), query, w, r)
	} else if authorize.ResponseType == "token" {
		// When response_type == "token", we will directly grant an access token
		authorizeResult, redirectURI, perr := router.manager.GrantAccessToken(authorize, accessPayload)
		if perr != nil && redirectURI != nil {
			routeUtils.ErrorRedirect(w, r, redirectURI, perr.PlutoCode, authorize.State, authorize.ResponseType)
			return nil
		} else if perr != nil {
			return perr
		}

		query := redirectURI.Query()
		query.Set("access_token", authorizeResult.AccessToken)
		query.Set("scopes", authorizeResult.Scopes)
		if authorize.State != "" {
			query.Set("state", authorizeResult.State)
		}

		routeUtils.RedirectWithFragment(authorize.RedirectURI, query, w, r)
	}

	return nil
}

func (router *Router) OAuthLogin(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	login := request.PasswordLogin{}

	if perr := routeUtils.GetRequestData(r, &login); perr != nil {
		return perr
	}

	var grantResult *manage.GrantResult

	var perr *perror.PlutoError
	if general.ValidMail(login.Account) {
		grantResult, perr = router.manager.MailPasswordLogin(login)
	} else {
		grantResult, perr = router.manager.NamePasswordLogin(login)
	}

	if perr != nil {
		return perr
	}

	cookie := &http.Cookie{
		Name:     "access_token",
		Value:    grantResult.AccessToken,
		Secure:   router.config.Server.CookieSecure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
		Expires:  time.Now().Add(60 * time.Minute),
	}
	http.SetCookie(w, cookie)

	loginRedirectURI := r.URL.Query().Get("login_redirect_uri")
	if loginRedirectURI == "" {
		loginRedirectURI = "/web/authorize"
	}

	loginRedirectURL, err := url.Parse(loginRedirectURI)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	var u *url.URL
	if loginRedirectURL.Hostname() == "" {
		u, err = url.Parse(loginRedirectURL.String())
		if err != nil {
			return perror.ServerError.Wrapper(err)
		}
	} else {
		u, err := url.Parse(routeUtils.GetBaseURL(r))
		if err != nil {
			return perror.ServerError.Wrapper(err)
		}
		u.Path = path.Join(u.Path, loginRedirectURL.Path)
	}

	redirectURL := fmt.Sprintf("%s%s", u.String(), routeUtils.GetQueryString(r.URL.Query(), "login_redirect_uri"))
	router.logger.Debug(redirectURL)

	routeUtils.RedirectWithQueryString(redirectURL, nil, w, r)

	return nil
}

func (router *Router) OAuthGetClient(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	accessPayload, perr := routeUtils.GetAccessPayload(r)

	if perr != nil {
		return perr
	}

	clients, perr := router.manager.OAuthGetClient(accessPayload)

	if perr != nil {
		return perr
	}

	data := make(map[string]interface{})

	data["clients"] = clients

	routeUtils.ResponseOK(data, w)

	return nil
}

func (router *Router) OAuthCreateClient(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	accessPayload, perr := routeUtils.GetAccessPayload(r)

	if perr != nil {
		return perr
	}

	occ := &request.OAuthCreateClient{}
	if perr := routeUtils.GetRequestData(r, occ); perr != nil {
		return perr
	}

	client, perr := router.manager.OAuthCreateClient(accessPayload, occ)

	if perr != nil {
		return perr
	}

	if perr := routeUtils.ResponseOK(client.Format(), w); perr != nil {
		return perr
	}

	return nil
}

func (router *Router) OAuthApproveClient(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	ocs := &request.OAuthClientStatus{}
	if perr := routeUtils.GetRequestData(r, ocs); perr != nil {
		return perr
	}

	client, perr := router.manager.UpdateOAuthClientStatus(ocs)

	if perr != nil {
		return perr
	}

	data := make(map[string]interface{})

	data["key"] = client.Key
	data["status"] = client.Status

	if perr := routeUtils.ResponseOK(data, w); perr != nil {
		return perr
	}

	return nil
}
