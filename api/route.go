package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AddRoute(router *httprouter.Router) {
	addUserRoute(router)
	addAuthRoute(router)
}

func addUserRoute(router *httprouter.Router) {
	router.POST("/api/user/register", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Register!\n")
	})

	router.POST("/api/user/login", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	})
}

func addAuthRoute(router *httprouter.Router) {
	router.POST("/api/auth/refresh", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	})
}
