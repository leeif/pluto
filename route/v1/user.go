package v1

import (
	"net/http"

	perror "github.com/leeif/pluto/datatype/pluto_error"

	"github.com/leeif/pluto/utils/mail"
	routeUtils "github.com/leeif/pluto/utils/route"

	"github.com/leeif/pluto/datatype/request"
)

func (router *Router) PasswordResetMail(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	rpm := request.ResetPasswordMail{}

	if err := routeUtils.GetRequestData(r, &rpm); err != nil {
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

		if err := ml.SendResetPassword(user.Mail.String, routeUtils.GetBaseURL(r), r.Header.Get("Accept-Language")); err != nil {
			router.logger.Error(err.LogError.Error())
		}
	}()

	routeUtils.ResponseOK(nil, w)

	return nil
}

func (router *Router) UserInfo(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	payload, perr := routeUtils.GetAccessPayload(r)
	if perr != nil {
		return perr
	}

	res, perr := router.manager.UserInfo(payload)

	if perr != nil {
		return perr
	}

	routeUtils.ResponseOK(res.Format(), w)

	return nil
}

func (router *Router) UpdateUserInfo(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	payload, perr := routeUtils.GetAccessPayload(r)
	if perr != nil {
		return perr
	}

	uui := request.UpdateUserInfo{}
	if err := routeUtils.GetRequestData(r, &uui); err != nil {
		return perr
	}

	err := router.manager.UpdateUserInfo(payload, uui)

	if err != nil {
		return err
	}

	routeUtils.ResponseOK(nil, w)

	return nil
}
