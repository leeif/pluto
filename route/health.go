package route

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (route *Route) healthCheckRoute(router *mux.Router) {
	router.Handle("/healthcheck", route.middleware.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		responseOK(nil, w)
	})).Methods("GET")
}
