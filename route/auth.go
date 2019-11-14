package route

import (
	"database/sql"
	"net/http"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/middleware"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/manage"
	"github.com/leeif/pluto/utils/rsa"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
)

func authRouter(router *mux.Router, db *sql.DB, config *config.Config, logger *log.PlutoLog) {
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
