package route

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/middleware"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/manage"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func userRouter(router *mux.Router, db *sql.DB, config *config.Config, logger *log.PlutoLog) {
	mw := middleware.NewMiddle(logger)
	manager := manage.NewManager(db, config, logger)

	router.Handle("/password/reset/mail", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rpm := request.ResetPasswordMail{}

		if err := getBody(r, &rpm); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		if err := manager.ResetPasswordMail(rpm, getBaseURL(r)); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		responseOK(nil, w)
	})).Methods("POST")

	router.Handle("/info/me", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		auth := strings.Fields(r.Header.Get("Authorization"))

		if len(auth) < 2 && strings.ToLower(auth[0]) != "jwt" {
			context.Set(r, "pluto_error", perror.InvalidJWTToekn)
			responseError(perror.InvalidJWTToekn, w)
			next(w, r)
			return
		}

		res, err := manager.UserInfo(auth[1])

		if err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		responseOK(formatUser(res), w)
	})).Methods("GET")

	router.Handle("/info/me/update", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		uui := request.UpdateUserInfo{}
		auth := strings.Fields(r.Header.Get("Authorization"))

		if err := getBody(r, &uui); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		if len(auth) < 2 && strings.ToLower(auth[0]) != "jwt" {
			context.Set(r, "pluto_error", perror.InvalidJWTToekn)
			responseError(perror.InvalidJWTToekn, w)
			next(w, r)
			return
		}

		err := manager.UpdateUserInfo(auth[1], uui)

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
