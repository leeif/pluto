package v1

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/manage"
	"github.com/leeif/pluto/utils/general"
	"github.com/leeif/pluto/utils/mail"
	routeUtils "github.com/leeif/pluto/utils/route"
)

func (router *Router) Login(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	login := request.PasswordLogin{}

	if err := routeUtils.GetRequestData(r, &login); err != nil {
		return err
	}

	var grantResult *manage.GrantResult

	var perr *perror.PlutoError
	if general.ValidMail(login.Account) {
		grantResult, perr = router.manager.MailPasswordLogin(login)
	} else {
		grantResult, perr = router.manager.NamePasswordLogin(login)
	}

	if perr != nil {
		return perr
	}

	routeUtils.ResponseOK(grantResult, w)

	return nil
}

func (router *Router) GoogleLoginMobile(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	login := request.GoogleMobileLogin{}

	if err := routeUtils.GetRequestData(r, &login); err != nil {
		return err
	}

	res, err := router.manager.GoogleLoginMobile(login)

	if err != nil {
		return err
	}

	routeUtils.ResponseOK(res, w)

	return nil
}

func (router *Router) WechatLoginMobile(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	login := request.WechatMobileLogin{}

	if err := routeUtils.GetRequestData(r, &login); err != nil {
		return err
	}

	res, err := router.manager.WechatLoginMobile(login)

	if err != nil {
		return err
	}

	routeUtils.ResponseOK(res, w)

	return nil
}

func (router *Router) AppleLoginMobile(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	login := request.AppleMobileLogin{}

	if err := routeUtils.GetRequestData(r, &login); err != nil {
		return err
	}

	res, err := router.manager.AppleLoginMobile(login)

	if err != nil {
		return err
	}

	routeUtils.ResponseOK(res, w)

	return nil
}

func (router *Router) BindMail(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	payload, perr := routeUtils.GetAccessPayload(r)
	if perr != nil {
		return perr
	}

	mb := &request.MailBinding{}

	if err := routeUtils.GetRequestData(r, &mb); err != nil {
		return err
	}

	if perr := router.manager.BindMail(mb, payload); perr != nil {
		return perr
	}

	return nil
}

func (router *Router) BindGoogle(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	payload, perr := routeUtils.GetAccessPayload(r)
	if perr != nil {
		return perr
	}

	gb := &request.GoogleBinding{}

	if err := routeUtils.GetRequestData(r, &gb); err != nil {
		return err
	}

	if perr := router.manager.BindGoogle(gb, payload); perr != nil {
		return perr
	}

	return nil
}

func (router *Router) BindApple(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	payload, perr := routeUtils.GetAccessPayload(r)
	if perr != nil {
		return perr
	}

	ab := &request.AppleBinding{}

	if err := routeUtils.GetRequestData(r, &ab); err != nil {
		return err
	}

	if perr := router.manager.BindApple(ab, payload); perr != nil {
		return perr
	}

	return nil
}

func (router *Router) BindWechat(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	payload, perr := routeUtils.GetAccessPayload(r)
	if perr != nil {
		return perr
	}

	wb := &request.WechatBinding{}

	if err := routeUtils.GetRequestData(r, &wb); err != nil {
		return err
	}

	if perr := router.manager.BindWechat(wb, payload); perr != nil {
		return perr
	}

	return nil
}

func (router *Router) Unbind(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	payload, perr := routeUtils.GetAccessPayload(r)
	if perr != nil {
		return perr
	}

	ub := &request.UnBinding{}

	if err := routeUtils.GetRequestData(r, &ub); err != nil {
		return err
	}

	if perr := router.manager.Unbind(ub, payload); perr != nil {
		return perr
	}

	return nil
}

func (router *Router) PasswordResetMail(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	rpm := request.ResetPasswordMail{}

	if err := routeUtils.GetRequestData(r, &rpm); err != nil {
		return err
	}

	perr := router.manager.ResetPasswordMail(rpm)
	if perr != nil {
		return perr
	}

	go func() {
		ml, err := mail.NewMail(router.config)
		if err != nil {
			router.logger.Error(err.LogError.Error())
		}

		if err := ml.SendResetPassword(rpm.Mail, routeUtils.GetBaseURL(r), r.Header.Get("Accept-Language")); err != nil {
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

	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userID"])

	if err != nil {
		return perror.BadRequest.Wrapper(err)
	}

	res, perr := router.manager.UserInfo(uint(userID), payload)

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

	vars := mux.Vars(r)

	userID, err := strconv.Atoi(vars["userID"])

	if err != nil {
		return perror.BadRequest.Wrapper(err)
	}

	if uint(userID) != payload.UserID {
		return perror.InvalidAccessToken
	}

	uui := request.UpdateUserInfo{}
	if err := routeUtils.GetRequestData(r, &uui); err != nil {
		return perr
	}

	if perr := router.manager.UpdateUserInfo(payload, uui); perr != nil {
		return perr
	}

	routeUtils.ResponseOK(nil, w)

	return nil
}

func (router *Router) Register(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	register := request.MailRegister{}

	if err := routeUtils.GetRequestData(r, &register); err != nil {
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
		if err := ml.SendRegisterVerify(user.ID, register.Mail, routeUtils.GetBaseURL(r), r.Header.Get("Accept-Language")); err != nil {
			router.logger.Error("send mail failed: " + err.LogError.Error())
		}
	}()

	routeUtils.ResponseOK(respBody, w)

	return nil
}

func (router *Router) VerifyMail(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	rvm := request.RegisterVerifyMail{}

	if perr := routeUtils.GetRequestData(r, &rvm); perr != nil {
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
		if err := ml.SendRegisterVerify(user.ID, rvm.Mail, routeUtils.GetBaseURL(r), r.Header.Get("Accept-Language")); err != nil {
			router.logger.Error(err.LogError.Error())
		}
	}()

	routeUtils.ResponseOK(nil, w)

	return nil
}
