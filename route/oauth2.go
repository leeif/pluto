package route

import (
	"net/http"
)

func (router *Router) Oauth2Tokens(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
}

func (router *Router) Oauth2Authorize(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
}

func (router *Router) Oauth2Login(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
}
