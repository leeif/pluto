package route

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/leeif/pluto/datatype/request"
)

func (router *Router) OAuth2Tokens(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	tokens := &request.OAuth2Tokens{}

	if err := getBody(r, tokens); err != nil {
		context.Set(r, "pluto_error", err)
		next(w, r)
		responseError(err, w)
		return
	}

	grantTypes := map[string]func(*request.OAuth2Tokens){
		"authorization_code": router.manager.AuthorizationCodeGrant,
		"password":           router.manager.PasswordGrant,
		"client_credentials": router.manager.ClientCredentialGrant,
		"refresh_token":      router.manager.RefreshTokenGrant,
	}

	grantHandler, ok := grantTypes[tokens.GrantType]
	if !ok {
		return
	}

	grantHandler(tokens)
}

func (router *Router) OAuth2Authorize(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authorize := &request.OAuth2Authorize{}

	if err := getBody(r, authorize); err != nil {
		context.Set(r, "pluto_error", err)
		next(w, r)
		responseError(err, w)
		return
	}

	if authorize.ResponseType == "code" {
		router.manager.GrantAuthorizationCode(authorize)
	} else if authorize.ResponseType == "token" {
		router.manager.GrantAccessToken(authorize)
	}
}

func (router *Router) OAuth2Login(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

}
