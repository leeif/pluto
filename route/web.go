package route

import (
	"net/http"

	"github.com/gorilla/mux"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
)

func (router *Router) registrationVerifyPage(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
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
		if err := responseHTMLError("register_verify_result.html", data, w, data.Error.HTTPCode); err != nil {
			router.logger.Error(err.Error())
		}
	} else {
		if err := responseHTMLOK("register_verify_result.html", data, w); err != nil {
			router.logger.Error(err.Error())
		}
	}

	return nil
}

func (router *Router) resetPasswordPage(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	vars := mux.Vars(r)
	token := vars["token"]

	type Data struct {
		Error *perror.PlutoError
		Token string
	}

	data := &Data{Token: token}
	if err := router.manager.ResetPasswordPage(token); err != nil {
		if err.PlutoCode == perror.ServerError.PlutoCode {
			return err
		}
		data.Error = err
		goto responseHTML
	}

responseHTML:
	if data.Error != nil {
		if err := responseHTMLError("password_reset.html", data, w, data.Error.HTTPCode); err != nil {
			router.logger.Error(err.Error())
		}
	} else if data.Error == nil {
		if err := responseHTMLOK("password_reset.html", data, w); err != nil {
			router.logger.Error(err.Error())
		}
	}

	return nil
}

func (router *Router) resetPassword(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	rpw := request.ResetPasswordWeb{}
	vars := mux.Vars(r)
	token := vars["token"]

	type Data struct {
		Error *perror.PlutoError
	}

	data := &Data{}

	if err := getBody(r, &rpw); err != nil {
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
		if err := responseHTMLError("password_reset_result.html", data, w, data.Error.HTTPCode); err != nil {
			router.logger.Error(err.Error())
		}
	} else {
		if err := responseHTMLOK("password_reset_result.html", data, w); err != nil {
			router.logger.Error(err.Error())
		}
	}

	return nil
}
