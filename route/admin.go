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
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
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
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
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
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
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
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
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
	application, err := router.manager.CreateApplication(ca)
	if err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	res := make(map[string]string)
	res["application_identifier"] = application.Identifier
	if err := responseOK(res, w); err != nil {
		router.logger.Error(err.Error())
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
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
	}
}

func (router *Router) ListApplications(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	applications, err := router.manager.ListApplications()

	if err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	apps := make([]map[string]interface{}, 0)

	for _, application := range applications {
		m := make(map[string]interface{})
		m["id"] = application.ID
		m["identifier"] = application.Identifier
		m["name"] = application.Name
		apps = append(apps, m)
	}

	if err := responseOK(apps, w); err != nil {
		router.logger.Error(err.Error())
	}
}

func (router *Router) ListRoles(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	lr := request.ListRoles{}
	// get paramter from url query
	if err := getQuery(r, &lr); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	roles, err := router.manager.ListRoles(lr.AppID)

	if err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	res := make([]map[string]interface{}, 0)

	for _, role := range roles {
		m := make(map[string]interface{})
		m["id"] = role.ID
		m["name"] = role.Name
		res = append(res, m)
	}

	if err := responseOK(res, w); err != nil {
		router.logger.Error(err.Error())
	}
}

func (router *Router) ListScopes(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ls := request.ListScopes{}
	// get paramter from url query
	if err := getQuery(r, &ls); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	scopes, err := router.manager.ListScopes(ls.AppID)

	if err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	res := make([]map[string]interface{}, 0)

	for _, scope := range scopes {
		m := make(map[string]interface{})
		m["id"] = scope.ID
		m["name"] = scope.Name
		res = append(res, m)
	}

	if err := responseOK(res, w); err != nil {
		router.logger.Error(err.Error())
	}
}

func (router *Router) UpdateUserRole(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	uur := request.UpdateUserRole{}

	if err := getBody(r, &uur); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	err := router.manager.UpdateUserRole(uur)

	if err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
	}
}

func (router *Router) FindUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fu := request.FindUser{}
	if err := getQuery(r, &fu); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
}
