package route

import (
	"net/http"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/middleware"

	b64 "encoding/base64"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/manage"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/mail"
	"github.com/leeif/pluto/utils/rsa"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func Router(router *mux.Router, db *gorm.DB, config *config.Config, logger *log.PlutoLog) {

	// user router
	user := router.PathPrefix("/api/user").Subrouter()
	userRouter(user, db, config, logger)

	// auth router
	auth := router.PathPrefix("/api/auth").Subrouter()
	authRouter(auth, db, config, logger)

	// web router
	web := router.PathPrefix("/").Subrouter()
	webRouter(web, db, config, logger)

	// health check router
	health := router.PathPrefix("/").Subrouter()
	healthCheckRouter(health, logger)
}

func userRouter(router *mux.Router, db *gorm.DB, config *config.Config, logger *log.PlutoLog) {
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
		go ml.SendRegisterVerify(userID, register.Mail, r.Host)
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

	router.Handle("/password/reset/mail", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rpm := request.ResetPasswordMail{}

		if err := getBody(r, &rpm); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		if err := manager.ResetPasswordMail(rpm, r.Host); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		responseOK(nil, w)
	})).Methods("POST")

	router.Handle("/password/reset", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rp := request.ResetPassword{}

		if err := getBody(r, &rp); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}
		if err := manager.ResetPassword(rp); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		// generate JWT for password reset result page
		token, err := jwt.GenerateJWT(jwt.Head{Type: jwt.PASSWORDRESETRESULT}, &jwt.PasswordResetResultPayload{Message: "Success"}, 10*60)
		if err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}
		res := make(map[string]string)
		res["redirect"] = "/password/reset/result/" + b64.StdEncoding.EncodeToString([]byte(token))

		responseOK(res, w)
	})).Methods("POST")

	router.Handle("/login", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		login := request.MailLogin{}

		if err := getBody(r, &login); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		res, err := manager.LoginWithEmail(login)

		if err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		responseOK(res, w)
	})).Methods("POST")

	router.Handle("/login/google", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		login := request.GoogleLogin{}

		if err := getBody(r, &login); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		res, err := manager.LoginWithGoogle(login)

		if err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		responseOK(res, w)
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

		responseOK(res, w)
	})).Methods("GET")
}

func authRouter(router *mux.Router, db *gorm.DB, config *config.Config, logger *log.PlutoLog) {
	mw := middleware.NewMiddle(logger)
	manager := manage.NewManager(db, config, logger)

	router.Handle("/refresh", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rat := request.RefreshAccessToken{}
		if err := getBody(r, &rat); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		res, err := manager.RefreshAccessToken(rat)

		if err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		responseOK(res, w)
	})).Methods("POST")

	router.Handle("/publickey", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		res := make(map[string]string)
		pbkey, err := rsa.GetPublicKey()

		if err != nil {
			perr := perror.ServerError.Wrapper(err)
			responseError(perr, w)
			next(w, r)
			return
		}

		res["public_key"] = pbkey
		responseOK(res, w)
	})).Methods("GET")
}
