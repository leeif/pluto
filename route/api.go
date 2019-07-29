package route

import (
	"fmt"
	"net/http"
	"os"

	"github.com/leeif/pluto/database"

	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/manage"

	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter()
	user := router.PathPrefix("/api/user").Subrouter()
	userRoute(user)

	auth := router.PathPrefix("/api/auth").Subrouter()
	authRoute(auth)

	return router
}

func userRoute(router *mux.Router) {
	db, err := database.GetDatabase()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	router.Handle("/register", negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			register := request.MailRegister{}

			if err := getBody(r, &register); err != nil {
				responseError(err, w)
			}

			if err := manage.RegisterWithEmail(db, register); err != nil {
				responseError(err, w)
			} else {
				respBody := make(map[string]interface{})
				respBody["mail"] = register.Mail
				responseOK(respBody, w)
			}
		}),
	)).Methods("POST")

	router.Handle("/login", negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			login := request.MailLogin{}
			getBody(r, &login)
			if jwtToken, err := manage.LoginWithEmail(db, login); err != nil {
				fmt.Println(err.Err.Error())
				responseError(err, w)
			} else {
				respBody := make(map[string]interface{})
				respBody["jwt"] = jwtToken
				responseOK(respBody, w)
			}
		}),
	)).Methods("POST")
}

func authRoute(router *mux.Router) {
	router.Handle("/refresh", negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			fmt.Println("refresh")
		}),
	))
}
