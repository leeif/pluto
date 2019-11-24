package route

import (
	"net/http"

	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/utils/mail"

	"github.com/gorilla/context"
)

func (router *Router) register(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	register := request.MailRegister{}

	if err := getBody(r, &register); err != nil {
		context.Set(r, "pluto_error", err)
		next(w, r)
		responseError(err, w)
		return
	}

	user, err := router.manager.RegisterWithEmail(register)
	if err != nil {
		// set err to context for log
		context.Set(r, "pluto_error", err)
		next(w, r)
		responseError(err, w)
		return
	}

	respBody := make(map[string]interface{})
	respBody["mail"] = register.Mail
	respBody["verified"] = user.Verified.Bool
	go func() {
		if router.config.Server.SkipRegisterVerifyMail {
			router.logger.Info("skip sending register mail")
			return
		}
		ml, err := mail.NewMail(router.config)
		if err != nil {
			router.logger.Error("send mail failed: " + err.LogError.Error())
		}
		if err := ml.SendRegisterVerify(user.ID, register.Mail, getBaseURL(r)); err != nil {
			router.logger.Error("send mail failed: " + err.LogError.Error())
		}
	}()
	responseOK(respBody, w)
}

func (router *Router) verifyMail(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rvm := request.RegisterVerifyMail{}

	if perr := getBody(r, &rvm); perr != nil {
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}

	user, perr := router.manager.RegisterVerifyMail(rvm)

	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}

	go func() {
		ml, err := mail.NewMail(router.config)
		if err != nil {
			router.logger.Error(err.LogError.Error())
		}
		if err := ml.SendRegisterVerify(user.ID, rvm.Mail, getBaseURL(r)); err != nil {
			router.logger.Error(err.LogError.Error())
		}
	}()
	responseOK(nil, w)
}
