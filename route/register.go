package route

import (
	"database/sql"
	"net/http"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/middleware"

	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/manage"
	"github.com/leeif/pluto/utils/mail"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func registerRouter(router *mux.Router, db *sql.DB, config *config.Config, logger *log.PlutoLog) {
	mw := middleware.NewMiddle(logger)
	manager := manage.NewManager(db, config, logger)

	router.Handle("/register", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		register := request.MailRegister{}

		if err := getBody(r, &register); err != nil {
			context.Set(r, "pluto_error", err)
			next(w, r)
			responseError(err, w)
			return
		}

		user, err := manager.RegisterWithEmail(register)
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
			if config.Server.SkipRegisterVerifyMail {
				logger.Info("skip sending register mail")
				return
			}
			ml, err := mail.NewMail(config)
			if err != nil {
				logger.Error(err.LogError.Error())
			}
			if err := ml.SendRegisterVerify(user.ID, register.Mail, getBaseURL(r)); err != nil {
				logger.Error(err.LogError.Error())
			}
		}()
		responseOK(respBody, w)
	})).Methods("POST")

	router.Handle("/register/verify/mail", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rvm := request.RegisterVerifyMail{}

		if err := getBody(r, &rvm); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		user, err := manager.RegisterVerifyMail(rvm)

		if err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		go func() {
			ml, err := mail.NewMail(config)
			if err != nil {
				logger.Error(err.LogError.Error())
			}
			if err := ml.SendRegisterVerify(user.ID, rvm.Mail, getBaseURL(r)); err != nil {
				logger.Error(err.LogError.Error())
			}
		}()

		responseOK(nil, w)
	})).Methods("POST")
}
