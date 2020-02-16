package route

import (
	"net/http"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
)

func (router *Router) login(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	login := request.MailLogin{}

	if err := getBody(r, &login); err != nil {
		return err
	}

	res, err := router.manager.EmailLogin(login)

	if err != nil {
		return err
	}

	responseOK(res, w)

	return nil
}

func (router *Router) googleLoginMobile(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	login := request.GoogleMobileLogin{}

	if err := getBody(r, &login); err != nil {
		return err
	}

	res, err := router.manager.GoogleLoginMobile(login)

	if err != nil {
		return err
	}

	responseOK(res, w)

	return nil
}

func (router *Router) wechatLoginMobile(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	login := request.WechatMobileLogin{}

	if err := getBody(r, &login); err != nil {
		return err
	}

	res, err := router.manager.WechatLoginMobile(login)

	if err != nil {
		return err
	}

	responseOK(res, w)

	return nil
}

func (router *Router) appleLoginMobile(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	login := request.AppleMobileLogin{}

	if err := getBody(r, &login); err != nil {
		return err
	}

	res, err := router.manager.AppleLoginMobile(login)

	if err != nil {
		return err
	}

	responseOK(res, w)

	return nil
}
