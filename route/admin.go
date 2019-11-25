package route

import (
	"net/http"

	"github.com/leeif/pluto/datatype/request"

	"github.com/gorilla/context"
)

func (router *Router) CreateRole(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	cr := request.CreateRole{}
	if err := getBody(r, &cr); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := router.manager.CreateRole(cr); err != nil {

	}
}

func (router *Router) CreateScope(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	cs := request.CreateScope{}
	if err := getBody(r, &cs); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := router.manager.CreateScope(cs); err != nil {

	}
}

func (router *Router) AttachScope(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rs := request.RoleScope{}
	if err := getBody(r, &rs); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := router.manager.AttachScope(rs); err != nil {

	}
}

func (router *Router) DetachScope(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rs := request.RoleScope{}
	if err := getBody(r, &rs); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := router.manager.DetachScope(rs); err != nil {

	}
}

func (router *Router) CreateApplication(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ca := request.CreateApplication{}
	if err := getBody(r, &ca); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := router.manager.CreateApplication(ca); err != nil {

	}
}
func (router *Router) AttachRole(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ar := request.ApplicationRole{}
	if err := getBody(r, &ar); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := router.manager.AttachRole(ar); err != nil {

	}
}

func (router *Router) ApplicationDefaultRole(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ar := request.ApplicationRole{}
	if err := getBody(r, &ar); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := router.manager.ApplicationDefaultRole(ar); err != nil {

	}
}
