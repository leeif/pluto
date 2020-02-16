package route

import (
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"net/http"

	"github.com/leeif/pluto/utils/mail"

	"github.com/leeif/pluto/datatype/request"
)

func (router *Router) passwordResetMail(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	rpm := request.ResetPasswordMail{}

	if err := getBody(r, &rpm); err != nil {
		return err
	}

	user, err := router.manager.ResetPasswordMail(rpm)
	if err != nil {
		return err
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

	return nil
}

func (router *Router) userInfo(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	payload, perr := getAccessPayload(r)
	if perr != nil {
		return perr
	}

	res, perr := router.manager.UserInfo(payload)

	if perr != nil {
		return perr
	}

	responseOK(res.Format(), w)

	return nil
}

func (router *Router) updateUserInfo(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	payload, perr := getAccessPayload(r)
	if perr != nil {
		return perr
	}

	uui := request.UpdateUserInfo{}
	if err := getBody(r, &uui); err != nil {
		return perr
	}

	err := router.manager.UpdateUserInfo(payload, uui)

	if err != nil {
		return err
	}

	responseOK(nil, w)

	return nil
}
