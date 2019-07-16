package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AddRoute(router *httprouter.Router) {
	router.POST("/api/register", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		fmt.Fprint(w, "Register!\n")
	})
}
