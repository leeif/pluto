package route

import (
	"fmt"
	"net/http"

	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

func GetAPIRouter(router *mux.Router) {
	router.PathPrefix("/api").Handler(userRoute())
	router.PathPrefix("/api").Handler(authRoute())
}

func userRoute() http.Handler {
	router := mux.NewRouter()
	router.Handle("/user/register", negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			fmt.Println("register")
		}),
	))
	router.Handle("/user/login", negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			fmt.Println("login")
		}),
	))
	return router
}

func authRoute() http.Handler {
	router := mux.NewRouter()
	router.Handle("/user/auth/refresh", negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			fmt.Println("refresh")
		}),
	))
	return router
}
