package route

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/leeif/pluto/datatype/request"
)

func (router *Router) login(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	login := request.MailLogin{}

	if err := getBody(r, &login); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	res, err := router.manager.EmailLogin(login)

	if err != nil {
		// set err to context for log
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	responseOK(res, w)
}

func (router *Router) googleLoginMobile(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	login := request.GoogleMobileLogin{}

	if err := getBody(r, &login); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	res, err := router.manager.GoogleLoginMobile(login)

	if err != nil {
		// set err to context for log
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	responseOK(res, w)
}

func (router *Router) wechatLoginMobile(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	login := request.WechatMobileLogin{}

	if err := getBody(r, &login); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	res, err := router.manager.WechatLoginMobile(login)

	if err != nil {
		// set err to context for log
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	responseOK(res, w)
}

func (router *Router) appleLoginMobile(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	login := request.AppleMobileLogin{}

	if err := getBody(r, &login); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	res, err := router.manager.AppleLoginMobile(login)

	if err != nil {
		// set err to context for log
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	responseOK(res, w)
}
