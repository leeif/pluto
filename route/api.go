package route

import (
	"fmt"
	"net/http"
	"os"

	b64 "encoding/base64"

	"github.com/go-kit/kit/log"
	"github.com/leeif/pluto/database"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/manage"
	"github.com/leeif/pluto/middleware"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/mail"
	"github.com/leeif/pluto/utils/rsa"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

type Route struct {
	Logger     log.Logger
	middleware middleware.Middleware
}

func (route *Route) GetRouter(logger log.Logger) *mux.Router {

	// new a middleware this a logger
	route.middleware = middleware.Middleware{Logger: route.Logger}

	router := mux.NewRouter()
	user := router.PathPrefix("/api/user").Subrouter()
	route.userRoute(user)

	auth := router.PathPrefix("/api/auth").Subrouter()
	route.authRoute(auth)

	web := router.PathPrefix("/").Subrouter()
	route.webRoute(web)

	return router
}

func (route *Route) userRoute(router *mux.Router) {
	db, err := database.GetDatabase()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	router.Handle("/register", route.middleware.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		register := request.MailRegister{}

		if err := getBody(r, &register); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		userID, err := manage.RegisterWithEmail(db, register)
		if err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		respBody := make(map[string]interface{})
		respBody["mail"] = register.Mail
		responseOK(respBody, w)
		mail.SendRegisterVerify(userID, register.Mail)
	})).Methods("POST")

	router.Handle("/register/verify/mail", route.middleware.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rvm := request.RegisterVerifyMail{}

		if err := getBody(r, &rvm); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		err := manage.RegisterVerifyMail(db, rvm)

		if err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		responseOK(nil, w)
	})).Methods("POST")

	router.Handle("/password/reset/mail", route.middleware.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rpm := request.ResetPasswordMail{}

		if err := getBody(r, &rpm); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		if err := manage.ResetPasswordMail(db, rpm); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		responseOK(nil, w)
	})).Methods("POST")

	router.Handle("/password/reset", route.middleware.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rp := request.ResetPassword{}

		if err := getBody(r, &rp); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}
		if err := manage.ResetPassword(db, rp); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		// generate JWT for password reset result page
		token, err := jwt.GenerateJWT(jwt.Head{Type: jwt.PASSWORDRESETRESULT}, &jwt.PasswordResetResultPayload{Message: "Success"}, 10*60)
		if err != nil {
			context.Set(r, "pluto_error", err)
			responseError(perror.NewServerError(err), w)
			next(w, r)
			return
		}
		res := make(map[string]string)
		res["redirect"] = "/password/reset/result/" + b64.StdEncoding.EncodeToString([]byte(token))

		responseOK(res, w)
	})).Methods("POST")

	router.Handle("/login", route.middleware.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		login := request.MailLogin{}

		if err := getBody(r, &login); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		res, err := manage.LoginWithEmail(db, login)

		if err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		responseOK(res, w)
	})).Methods("POST")
}

func (route *Route) authRoute(router *mux.Router) {
	db, err := database.GetDatabase()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	router.Handle("/refresh", route.middleware.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rat := request.RefreshAccessToken{}
		if err := getBody(r, &rat); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		res, err := manage.RefreshAccessToken(db, rat)

		if err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}

		responseOK(res, w)
	})).Methods("POST")

	router.Handle("/publickey", route.middleware.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		res := make(map[string]string)
		pbkey, err := rsa.GetPublicKey()

		if err != nil {
			perr := perror.NewServerError(err)
			responseError(perr, w)
			next(w, r)
			return
		}

		res["public_key"] = pbkey
		responseOK(res, w)
	})).Methods("GET")
}
