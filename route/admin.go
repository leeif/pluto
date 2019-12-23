package route

import (
	"net/http"

	"github.com/leeif/pluto/datatype/request"

	"github.com/gorilla/context"
)

func (router *Router) CreateRole(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	_, perr := getAccessPayload(r)
	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}
	cr := request.CreateRole{}
	if err := getBody(r, &cr); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if _, err := router.manager.CreateRole(cr); err != nil {
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
	_, perr := getAccessPayload(r)
	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}
	cs := request.CreateScope{}
	if err := getBody(r, &cs); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if _, err := router.manager.CreateScope(cs); err != nil {
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
	_, perr := getAccessPayload(r)
	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}
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

func (router *Router) RoleScopeBatchUpdate(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	_, perr := getAccessPayload(r)
	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}
	rsbu := request.RoleScopeBatchUpdate{}
	if err := getBody(r, &rsbu); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := router.manager.RoleScopeBatchUpdate(rsbu); err != nil {
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
	_, perr := getAccessPayload(r)
	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}
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
	_, perr := getAccessPayload(r)
	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}
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
	res := make(map[string]interface{})
	res["application"] = application
	if err := responseOK(res, w); err != nil {
		router.logger.Error(err.Error())
	}
}

func (router *Router) ApplicationDefaultRole(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	_, perr := getAccessPayload(r)
	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}
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

func (router *Router) RoleDefaultScope(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	_, perr := getAccessPayload(r)
	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}
	rs := request.RoleScope{}
	if err := getBody(r, &rs); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	if err := router.manager.RoleDefaultScope(rs); err != nil {
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
		m["name"] = application.Name
		m["default_role"] = application.DefaultRole
		apps = append(apps, m)
	}

	res := make(map[string]interface{})
	res["applications"] = apps

	if err := responseOK(res, w); err != nil {
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

	er, err := router.manager.ListRoles(lr.AppID)

	if err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	if err := responseOK(er.Format(), w); err != nil {
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

	es, err := router.manager.ListScopes(ls.AppID)
	if err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	s := make([]map[string]interface{}, 0)

	app := make(map[string]interface{})
	app["id"] = es.Application.ID
	app["name"] = es.Application.Name
	app["default_role"] = es.Application.DefaultRole

	for _, scope := range es.Scopes {
		m := make(map[string]interface{})
		m["id"] = scope.ID
		m["name"] = scope.Name
		s = append(s, m)
	}

	res := make(map[string]interface{})
	res["application"] = app
	res["scopes"] = s

	if err := responseOK(res, w); err != nil {
		router.logger.Error(err.Error())
	}
}

func (router *Router) SetUserRole(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	_, perr := getAccessPayload(r)
	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}

	ur := request.UserRole{}
	if err := getBody(r, &ur); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	err := router.manager.SetUserRole(ur)

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
	_, perr := getAccessPayload(r)
	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}
	fu := &request.FindUser{}
	if err := getQuery(r, fu); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	founds, err := router.manager.FindUser(fu)
	if err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
	users := make([]map[string]interface{}, 0)
	for _, found := range founds {
		users = append(users, found.Format())
	}
	res := make(map[string]interface{})
	res["users"] = users
	if err := responseOK(res, w); err != nil {
		router.logger.Error(err.Error())
	}
}

func (router *Router) UsersCount(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	_, perr := getAccessPayload(r)
	if perr != nil {
		// set err to context for log
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}
	res, perr := router.manager.UsersCount()
	if perr != nil {
		context.Set(r, "pluto_error", perr)
		responseError(perr, w)
		next(w, r)
		return
	}
	if err := responseOK(res, w); err != nil {
		router.logger.Error(err.Error())
	}
}
