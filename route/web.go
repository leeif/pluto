package route

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/leeif/pluto/config"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/manage"
	"github.com/leeif/pluto/middleware"
)

func webRouter(router *mux.Router, db *sql.DB, config *config.Config, logger *log.PlutoLog) {

	mw := middleware.NewMiddle(logger)
	manager := manage.NewManager(db, config, logger)

	router.Handle("/mail/verify/{token}", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		vars := mux.Vars(r)
		token := vars["token"]

		type Data struct {
			Error *perror.PlutoError
		}
		data := &Data{}

		if err := manager.RegisterVerify(token); err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			next(w, r)
			data.Error = err
			goto responseHTML
		}

	responseHTML:

		if data.Error != nil && data.Error.PlutoCode == perror.ServerError.PlutoCode {
			if err := responseHTMLError("error.html", nil, w, data.Error.HTTPCode); err != nil {
				logger.Error(err.Error())
			}
		} else if data.Error != nil {
			if err := responseHTMLError("register_verify_result.html", data, w, data.Error.HTTPCode); err != nil {
				logger.Error(err.Error())
			}
		} else {
			if err := responseHTMLOK("register_verify_result.html", data, w); err != nil {
				logger.Error(err.Error())
			}
		}

	})).Methods("GET")

	router.Handle("/password/reset/{token}", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		vars := mux.Vars(r)
		token := vars["token"]

		type Data struct {
			Error *perror.PlutoError
			Token string
		}

		data := &Data{Token: token}
		if err := manager.ResetPasswordPage(token); err != nil {
			context.Set(r, "pluto_error", err)
			next(w, r)
			data.Error = err
			goto responseHTML
		}

	responseHTML:
		if data.Error != nil && data.Error.PlutoCode == perror.ServerError.PlutoCode {
			if err := responseHTMLError("error.html", nil, w, data.Error.HTTPCode); err != nil {
				logger.Error(err.Error())
			}
		} else if data.Error != nil {
			if err := responseHTMLError("password_reset.html", data, w, data.Error.HTTPCode); err != nil {
				logger.Error(err.Error())
			}
		} else if data.Error == nil {
			if err := responseHTMLOK("password_reset.html", data, w); err != nil {
				logger.Error(err.Error())
			}
		}

	})).Methods("GET")

	router.Handle("/password/reset/result", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rp := request.ResetPassword{}

		type Data struct {
			Error *perror.PlutoError
		}

		data := &Data{}

		if err := getBody(r, &rp); err != nil {
			context.Set(r, "pluto_error", err)
			next(w, r)
			data.Error = err
			goto responseHTML
		}

		if err := manager.ResetPassword(rp); err != nil {
			context.Set(r, "pluto_error", err)
			next(w, r)
			data.Error = err
			goto responseHTML
		}

	responseHTML:
		if data.Error != nil && data.Error.PlutoCode == perror.ServerError.PlutoCode {
			if err := responseHTMLError("error.html", nil, w, data.Error.HTTPCode); err != nil {
				logger.Error(err.Error())
			}
		} else if data.Error != nil {
			if err := responseHTMLError("password_reset_result.html", data, w, data.Error.HTTPCode); err != nil {
				logger.Error(err.Error())
			}
		} else {
			if err := responseHTMLOK("password_reset_result.html", data, w); err != nil {
				logger.Error(err.Error())
			}
		}

	})).Methods("POST")
}
