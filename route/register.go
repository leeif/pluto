package route

import (
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"net/http"

	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/utils/mail"
)

func (router *Router) register(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	register := request.MailRegister{}

	if err := getBody(r, &register); err != nil {
		return err
	}

	user, err := router.manager.RegisterWithEmail(register)
	if err != nil {
		return err
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

	return nil
}

func (router *Router) verifyMail(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	rvm := request.RegisterVerifyMail{}

	if perr := getBody(r, &rvm); perr != nil {
		return perr
	}

	user, perr := router.manager.RegisterVerifyMail(rvm)

	if perr != nil {
		return perr
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

	return nil
}
