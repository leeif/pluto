package route

import (
	"net/http"

	"github.com/leeif/pluto/utils/mail"

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
	payload, perr := getAccessPayload(r)
	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}

	res, perr := router.manager.UserInfo(payload)

	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}

	responseOK(res.Format(), w)
}

func (router *Router) updateUserInfo(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	payload, perr := getAccessPayload(r)
	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}

	uui := request.UpdateUserInfo{}
	if err := getBody(r, &uui); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	err := router.manager.UpdateUserInfo(payload, uui)

	if err != nil {
		// set err to context for log
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	responseOK(nil, w)
}
