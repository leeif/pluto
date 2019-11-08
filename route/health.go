package route

import (
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/middleware"
)

func healthCheckRouter(router *mux.Router, db *sql.DB, config *config.Config, logger *log.PlutoLog) {
	mw := middleware.NewMiddle(logger)
	router.Handle("/healthcheck", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		responseOK(nil, w)
	})).Methods("GET")
}
