package route

import (
	"net/http"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
)

func (router *Router) OAuth2Tokens(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	tokens := &request.OAuth2Tokens{}

	if err := getBody(r, tokens); err != nil {
		return err
	}

	grantTypes := map[string]func(*request.OAuth2Tokens){
		"authorization_code": router.manager.AuthorizationCodeGrant,
		"password":           router.manager.PasswordGrant,
		"client_credentials": router.manager.ClientCredentialGrant,
		"refresh_token":      router.manager.RefreshTokenGrant,
	}

	grantHandler, ok := grantTypes[tokens.GrantType]
	if !ok {
		// TODO
		return nil
	}

	grantHandler(tokens)

	return nil
}

func (router *Router) OAuth2Authorize(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	authorize := &request.OAuth2Authorize{}

	if err := getBody(r, authorize); err != nil {
		return nil
	}

	if authorize.RedirectURI == "" {

	}

	if authorize.ResponseType == "code" {
		// Create a new authorization code
		query, err := router.manager.GrantAuthorizationCode(authorize)
		if err != nil {
			return err
		}
		redirectWithQueryString(authorize.RedirectURI, query.Query(), w, r)
	} else if authorize.ResponseType == "token" {
		// When response_type == "token", we will directly grant an access token
		query, err := router.manager.GrantAccessToken(authorize)
		if err != nil {
			return err
		}
		redirectWithFragment(authorize.RedirectURI, query.Query(), w, r)
	}

	return nil
}

func (router *Router) OAuth2Login(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	return nil
}
