package route

import (
	"net/http"

	"github.com/jinzhu/gorm"
	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/middleware"

	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/manage"
	"github.com/leeif/pluto/utils/mail"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func registerRouter(router *mux.Router, db *gorm.DB, config *config.Config, logger *log.PlutoLog) {
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

		userID, err := manager.RegisterWithEmail(register)
		if err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			next(w, r)
			responseError(err, w)
			return
		}

		respBody := make(map[string]interface{})
		respBody["mail"] = register.Mail
		ml := mail.NewMail(config)
		go func() {
			if err := ml.SendRegisterVerify(userID, register.Mail, r.Host); err != nil {
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

		err := manager.RegisterVerifyMail(db, rvm, r.Host)

		if err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		responseOK(nil, w)
	})).Methods("POST")
}
