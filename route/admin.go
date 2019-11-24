package route

import (
	"database/sql"
	"net/http"

	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/manage"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/middleware"
)

func adminRouter(router *mux.Router, db *sql.DB, config *config.Config, logger *log.PlutoLog) {
	mw := middleware.NewMiddle(logger)
	manager := manage.NewManager(db, config, logger)
	router.Handle("/rbac/role/create", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		cr := request.CreateRole{}
		if err := getBody(r, &cr); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}
		if err := manager.CreateRole(cr); err != nil {

		}
	})).Methods("POST")

	router.Handle("/rbac/scope/create", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		cs := request.CreateScope{}
		if err := getBody(r, &cs); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}
		if err := manager.CreateScope(cs); err != nil {

		}
	})).Methods("POST")

	router.Handle("/rbac/scope/attach", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rs := request.RoleScope{}
		if err := getBody(r, &rs); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}
		if err := manager.AttachScope(rs); err != nil {

		}
	})).Methods("POST")

	router.Handle("/rbac/scope/detach", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		rs := request.RoleScope{}
		if err := getBody(r, &rs); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}
		if err := manager.DetachScope(rs); err != nil {

		}
	})).Methods("POST")

	router.Handle("/application/create", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ca := request.CreateApplication{}
		if err := getBody(r, &ca); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}
		if err := manager.CreateApplication(ca); err != nil {

		}
	})).Methods("POST")

	router.Handle("/application/role/attach", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ar := request.ApplicationRole{}
		if err := getBody(r, &ar); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}
		if err := manager.AttachRole(ar); err != nil {

		}
	})).Methods("POST")

	router.Handle("/application/role/default", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ar := request.ApplicationRole{}
		if err := getBody(r, &ar); err != nil {
			context.Set(r, "pluto_error", err)
			responseError(err, w)
			next(w, r)
			return
		}
		if err := manager.ApplicationDefaultRole(ar); err != nil {

		}
	})).Methods("POST")
}
