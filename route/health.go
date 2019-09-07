package route

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/middleware"
)

func healthCheckRouter(router *mux.Router, logger *log.PlutoLog) {
	mw := middleware.NewMiddle(logger)
	router.Handle("/healthcheck", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		responseOK(nil, w)
	})).Methods("GET")
}
