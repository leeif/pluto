package v1

import (
	"net/http"

	"github.com/gorilla/mux"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/utils/general"
	routeUtils "github.com/leeif/pluto/utils/route"
)

func (router *Router) RegistrationVerifyPage(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	vars := mux.Vars(r)
	token := vars["token"]

	type Data struct {
		Error *perror.PlutoError
	}
	data := &Data{}

	if err := router.manager.RegisterVerify(token); err != nil {
		if err.PlutoCode == perror.ServerError.PlutoCode {
			return err
		}
		data.Error = err
		goto responseHTML
	}

responseHTML:

	if data.Error != nil {
		if err := routeUtils.ResponseHTMLError("register_verify_result.html", data, r, w, data.Error.HTTPCode); err != nil {
			return err
		}
	} else {
		if err := routeUtils.ResponseHTMLOK("register_verify_result.html", data, r, w); err != nil {
			return err
		}
	}

	return nil
}

func (router *Router) ResetPasswordPage(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	vars := mux.Vars(r)
	token := vars["token"]

	type Data struct {
		Error   *perror.PlutoError
		BaseURL string
		Token   string
	}

	data := &Data{
		Token:   token,
		BaseURL: routeUtils.GetBaseURL(r),
	}

	if err := router.manager.ResetPasswordPage(token); err != nil {
		if err.PlutoCode == perror.ServerError.PlutoCode {
			return err
		}
		data.Error = err
		goto responseHTML
	}

responseHTML:
	if data.Error != nil {
		if err := routeUtils.ResponseHTMLError("password_reset.html", data, r, w, data.Error.HTTPCode); err != nil {
			return err
		}
	} else if data.Error == nil {
		if err := routeUtils.ResponseHTMLOK("password_reset.html", data, r, w); err != nil {
			return err
		}
	}

	return nil
}

func (router *Router) ResetPassword(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	rpw := request.ResetPasswordWeb{}
	vars := mux.Vars(r)
	token := vars["token"]

	type Data struct {
		Error *perror.PlutoError
	}

	data := &Data{}

	if err := routeUtils.GetRequestData(r, &rpw); err != nil {
		if err.PlutoCode == perror.ServerError.PlutoCode {
			return err
		}
		data.Error = err
		goto responseHTML
	}

	if err := router.manager.ResetPassword(token, rpw); err != nil {
		if err.PlutoCode == perror.ServerError.PlutoCode {
			return err
		}
		data.Error = err
		goto responseHTML
	}

responseHTML:
	if data.Error != nil {
		if err := routeUtils.ResponseHTMLError("password_reset_result.html", data, r, w, data.Error.HTTPCode); err != nil {
			return err
		}
	} else {
		if err := routeUtils.ResponseHTMLOK("password_reset_result.html", data, r, w); err != nil {
			return err
		}
	}

	return nil
}

func (router *Router) LoginPage(w http.ResponseWriter, r *http.Request) *perror.PlutoError {

	query := r.URL.Query()

	if query.Get("app_id") == "" {
		query.Set("app_id", general.PlutoAdminApplication)
		routeUtils.RedirectWithQueryString(r.URL.Path, query, w, r)
		return nil
	}

	type Data struct {
		BaseURL     string
		Application string
	}

	data := &Data{
		BaseURL:     routeUtils.GetBaseURL(r),
		Application: query.Get("app_id"),
	}

	if err := routeUtils.ResponseHTMLOK("login.html", data, r, w); err != nil {
		return err
	}

	return nil
}

func (router *Router) AuthorizePage(w http.ResponseWriter, r *http.Request) *perror.PlutoError {

	query := r.URL.Query()

	accessPayload, perr := routeUtils.GetAccessPayload(r)
	if perr != nil {
		return perr
	}

	if accessPayload.AppID != query.Get("app_id") {
		query.Set("app_id", accessPayload.AppID)
		routeUtils.RedirectWithQueryString(r.URL.Path, query, w, r)
		return nil
	}

	type Data struct {
		AppID    string
		ClientID string
		BaseURL  string
		Token    bool
	}

	data := &Data{
		AppID:   accessPayload.AppID,
		BaseURL: routeUtils.GetBaseURL(r),
		Token:   query.Get("response_type") == "token",
	}

	clientID := query.Get("client_id")

	if clientID == "" {
		return perror.OAuthClientIDOrSecretNotFound
	}

	if _, perr := router.manager.GetClientByKey(clientID); perr != nil {
		return perr
	}

	data.ClientID = clientID

	routeUtils.ResponseHTMLOK("authorize.html", data, r, w)

	return nil
}
