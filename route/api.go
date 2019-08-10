package route

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-kit/kit/log"
	"github.com/leeif/pluto/database"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/manage"
	"github.com/leeif/pluto/middleware"
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
			responseError(err, w)
			return
		}

		if err := manage.RegisterWithEmail(db, register); err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			responseError(err, w)
		} else {
			respBody := make(map[string]interface{})
			respBody["mail"] = register.Mail
			responseOK(respBody, w)
		}
		next(w, r)
	})).Methods("POST")

	router.Handle("/login", route.middleware.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		login := request.MailLogin{}
		if err := getBody(r, &login); err != nil {
			responseError(err, w)
			return
		}
		if res, err := manage.LoginWithEmail(db, login); err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			responseError(err, w)
		} else {
			responseOK(res, w)
		}
		next(w, r)
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
			responseError(err, w)
			return
		}
		if res, err := manage.RefreshAccessToken(db, rat); err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			responseError(err, w)
		} else {
			responseOK(res, w)
		}
		next(w, r)
	})).Methods("POST")

	router.Handle("/publickey", route.middleware.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		res := make(map[string]string)
		pbkey, err := rsa.GetPublicKey()
		if err != nil {
			perr := perror.NewServerError(err)
			responseError(perr, w)
		} else {
			res["public_key"] = pbkey
			responseOK(res, w)
		}
		next(w, r)
	})).Methods("GET")
}
