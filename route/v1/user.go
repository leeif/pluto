package v1

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	perror "github.com/MuShare/pluto/datatype/pluto_error"
	"github.com/MuShare/pluto/datatype/request"
	"github.com/MuShare/pluto/manage"
	"github.com/MuShare/pluto/utils/general"
	"github.com/MuShare/pluto/utils/mail"
	routeUtils "github.com/MuShare/pluto/utils/route"
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

func (router *Router) WechatLoginWeb(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	query := r.URL.Query()

	code := query.Get("code")
	stateStr := query.Get("state")

	decoded, derr := base64.StdEncoding.DecodeString(stateStr)

	if derr != nil {
		return perror.BadRequest.Wrapper(fmt.Errorf("decode base64 stateStr error: %s", derr))
	}

	state := request.WechatWebLoginState{}
	if err := json.Unmarshal(decoded, &state); err != nil {
		return perror.BadRequest.Wrapper(fmt.Errorf("unmarshal state error:%s", err))
	}

	res, err := router.manager.WechatLoginWeb(state.AppID, code)

	if err != nil {
		return err
	}

	q := url.Values{}
	q.Add("token", res.AccessToken)
	http.Redirect(w, r, state.RedirectURL+"?"+q.Encode(), 302)

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

func (router *Router) Binding(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	payload, perr := routeUtils.GetAccessPayload(r)
	if perr != nil {
		return perr
	}

	binding := &request.Binding{}

	if err := routeUtils.GetRequestData(r, binding); err != nil {
		return err
	}

	switch binding.Type {
	case manage.MAILLOGIN:
		if binding.Mail == "" {
			return perror.BadRequest
		}
		perr = router.manager.BindMail(binding, payload)
		if perr == nil {
			if router.config.Server.SkipRegisterVerifyMail {
				router.logger.Info("skip sending register mail")
				return nil
			}
			ml, err := mail.NewMail(router.config, router.bundle)
			if err != nil {
				router.logger.Error("send mail failed: " + err.LogError.Error())
				return err
			}
			language := r.Header.Get("Accept-Language")
			appI18nName, err := router.manager.ApplicationI18nName(payload.AppID, language)
			if err != nil {
				router.logger.Error(err.LogError.Error())
				return err
			}
			if err := ml.SendRegisterVerify(payload.UserID, binding.Mail, routeUtils.GetBaseURL(r), language, appI18nName); err != nil {
				router.logger.Error("send mail failed: " + err.LogError.Error())
				return perror.SendMailFailure
			}
		}
	case manage.GOOGLELOGIN:
		if binding.IDToken == "" {
			return perror.BadRequest
		}
		perr = router.manager.BindGoogle(binding, payload)
	case manage.APPLELOGIN:
		if binding.Code == "" {
			return perror.BadRequest
		}
		perr = router.manager.BindApple(binding, payload)
	case manage.WECHATLOGIN:
		if binding.Code == "" {
			return perror.BadRequest
		}
		perr = router.manager.BindWechat(binding, payload)
	default:
		return perror.Forbidden
	}

	if perr != nil {
		return perr
	}

	routeUtils.ResponseOK(nil, w)

	return nil
}

func (router *Router) Unbinding(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	payload, perr := routeUtils.GetAccessPayload(r)
	if perr != nil {
		return perr
	}

	ub := &request.UnBinding{}

	if err := routeUtils.GetRequestData(r, ub); err != nil {
		return err
	}

	switch ub.Type {
	case manage.MAILLOGIN, manage.GOOGLELOGIN, manage.APPLELOGIN, manage.WECHATLOGIN:
		perr = router.manager.Unbind(ub, payload)
	default:
		return perror.Forbidden
	}

	if perr != nil {
		return perr
	}

	routeUtils.ResponseOK(nil, w)

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

	ml, err := mail.NewMail(router.config, router.bundle)
	if err != nil {
		router.logger.Error(err.LogError.Error())
		return err
	}
	language := r.Header.Get("Accept-Language")
	appI18nName, err := router.manager.ApplicationI18nName(rpm.AppName, language)
	if err != nil {
		router.logger.Error(err.LogError.Error())
		return err
	}
	if err := ml.SendResetPassword(rpm.Mail, routeUtils.GetBaseURL(r), language, appI18nName); err != nil {
		router.logger.Error(err.LogError.Error())
		return perror.SendMailFailure
	}

	routeUtils.ResponseOK(nil, w)

	return nil
}

func (router *Router) UserInfo(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	payload, perr := routeUtils.GetAccessPayload(r)
	if perr != nil {
		return perr
	}

	res, perr := router.manager.UserInfo(payload.UserID, payload)

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
		return err
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

	user, err := router.manager.RegisterWithEmail(register, false)
	if err != nil {
		return err
	}

	respBody := make(map[string]interface{})
	respBody["mail"] = register.Mail
	respBody["verified"] = user.Verified.Bool

	if router.config.Server.SkipRegisterVerifyMail {
		router.logger.Info("skip sending register mail")
		return nil
	}
	ml, err := mail.NewMail(router.config, router.bundle)
	if err != nil {
		router.logger.Error("send mail failed: " + err.LogError.Error())
		return err
	}
	language := r.Header.Get("Accept-Language")
	appI18nName, err := router.manager.ApplicationI18nName(register.AppName, language)
	if err != nil {
		router.logger.Error("send mail failed: " + err.LogError.Error())
		return err
	}
	if err := ml.SendRegisterVerify(user.ID, register.Mail, routeUtils.GetBaseURL(r), language, appI18nName); err != nil {
		router.logger.Error("send mail failed: " + err.LogError.Error())
		return perror.SendMailFailure
	}

	routeUtils.ResponseOK(respBody, w)

	return nil
}

func (router *Router) VerifyMail(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	rvm := request.RegisterVerifyMail{}

	if perr := routeUtils.GetRequestData(r, &rvm); perr != nil {
		return perr
	}

	binding, perr := router.manager.RegisterVerifyMail(rvm)

	if perr != nil {
		return perr
	}

	ml, err := mail.NewMail(router.config, router.bundle)
	if err != nil {
		router.logger.Error(err.LogError.Error())
		return err
	}
	language := r.Header.Get("Accept-Language")
	appI18nName, err := router.manager.ApplicationI18nName(rvm.AppName, language)
	if err != nil {
		router.logger.Error("send mail failed: " + err.LogError.Error())
		return err
	}
	if err := ml.SendRegisterVerify(binding.UserID, binding.Mail, routeUtils.GetBaseURL(r), language, appI18nName); err != nil {
		router.logger.Error(err.LogError.Error())
		return perror.SendMailFailure
	}

	routeUtils.ResponseOK(nil, w)

	return nil
}

func (router *Router) PublicUserInfo(w http.ResponseWriter, r *http.Request) *perror.PlutoError {

	pui := request.PublicUserInfos{}

	if perr := routeUtils.GetRequestData(r, &pui); perr != nil {
		return perr
	}

	res := make(map[string]map[string]interface{})
	for _, id := range pui.IDs {
		info, perr := router.manager.PublicUserInfo(id)

		if perr != nil {
			continue
		}

		res[id] = info
	}

	routeUtils.ResponseOK(res, w)

	return nil
}
