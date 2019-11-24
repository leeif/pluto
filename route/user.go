package route

import (
	"net/http"
	"strings"

	"github.com/leeif/pluto/utils/mail"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"

	"github.com/gorilla/context"
)

func (router *Router) passwordResetMail(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rpm := request.ResetPasswordMail{}

	if err := getBody(r, &rpm); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	user, err := router.manager.ResetPasswordMail(rpm)
	if err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	go func() {
		ml, err := mail.NewMail(router.config)
		if err != nil {
			router.logger.Error(err.LogError.Error())
		}

		if err := ml.SendResetPassword(user.Mail.String, getBaseURL(r)); err != nil {
			router.logger.Error(err.LogError.Error())
		}
	}()

	responseOK(nil, w)
}

func (router *Router) userInfo(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	auth := strings.Fields(r.Header.Get("Authorization"))

	if len(auth) < 2 && strings.ToLower(auth[0]) != "jwt" {
		context.Set(r, "pluto_error", perror.InvalidJWTToekn)
		responseError(perror.InvalidJWTToekn, w)
		next(w, r)
		return
	}

	res, err := router.manager.UserInfo(auth[1])

	if err != nil {
		// set err to context for log
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	responseOK(formatUser(res), w)
}

func (router *Router) updateUserInfo(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	uui := request.UpdateUserInfo{}
	auth := strings.Fields(r.Header.Get("Authorization"))

	if err := getBody(r, &uui); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	if len(auth) < 2 && strings.ToLower(auth[0]) != "jwt" {
		context.Set(r, "pluto_error", perror.InvalidJWTToekn)
		responseError(perror.InvalidJWTToekn, w)
		next(w, r)
		return
	}

	err := router.manager.UpdateUserInfo(auth[1], uui)

	if err != nil {
		// set err to context for log
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	responseOK(nil, w)
}
