package route

import (
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"net/http"

	"github.com/leeif/pluto/datatype/request"
)

func (router *Router) CreateRole(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	_, perr := getAccessPayload(r)
	if perr != nil {
		return perr
	}
	cr := request.CreateRole{}
	if err := getBody(r, &cr); err != nil {
		return perr
	}
	if _, err := router.manager.CreateRole(cr); err != nil {
		return perr
	}
	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
	}
	return nil
}

func (router *Router) CreateScope(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	_, perr := getAccessPayload(r)
	if perr != nil {
		return perr
	}
	cs := request.CreateScope{}
	if err := getBody(r, &cs); err != nil {
		return perr
	}
	if _, err := router.manager.CreateScope(cs); err != nil {
		return perr
	}
	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
	}

	return perr
}

func (router *Router) AttachScope(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	_, perr := getAccessPayload(r)
	if perr != nil {
		return perr
	}
	rs := request.RoleScope{}
	if err := getBody(r, &rs); err != nil {
		return perr
	}
	if perr := router.manager.AttachScope(rs); perr != nil {
		return perr
	}
	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
	}

	return nil
}

func (router *Router) RoleScopeBatchUpdate(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	_, perr := getAccessPayload(r)
	if perr != nil {
		return perr
	}
	rsbu := request.RoleScopeBatchUpdate{}
	if err := getBody(r, &rsbu); err != nil {
		return perr
	}
	if err := router.manager.RoleScopeBatchUpdate(rsbu); err != nil {
		return perr
	}
	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
	}

	return nil
}

func (router *Router) DetachScope(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	_, perr := getAccessPayload(r)
	if perr != nil {
		return perr
	}
	rs := request.RoleScope{}
	if err := getBody(r, &rs); err != nil {
		return perr
	}
	if err := router.manager.DetachScope(rs); err != nil {
		return perr
	}
	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
	}

	return nil
}

func (router *Router) CreateApplication(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	_, perr := getAccessPayload(r)
	if perr != nil {
		return perr
	}
	ca := request.CreateApplication{}
	if perr := getBody(r, &ca); perr != nil {
		return perr
	}
	application, err := router.manager.CreateApplication(ca)
	if err != nil {
		return perr
	}
	res := make(map[string]interface{})
	res["application"] = application
	if err := responseOK(res, w); err != nil {
		router.logger.Error(err.Error())
	}

	return nil
}

func (router *Router) ApplicationDefaultRole(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	_, perr := getAccessPayload(r)
	if perr != nil {
		return perr
	}
	ar := request.ApplicationRole{}
	if err := getBody(r, &ar); err != nil {
		return err
	}
	if err := router.manager.ApplicationDefaultRole(ar); err != nil {
		return err
	}
	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
	}

	return nil
}

func (router *Router) RoleDefaultScope(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	_, perr := getAccessPayload(r)
	if perr != nil {
		return perr
	}
	rs := request.RoleScope{}
	if err := getBody(r, &rs); err != nil {
		return err
	}
	if err := router.manager.RoleDefaultScope(rs); err != nil {
		return err
	}
	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
	}

	return nil
}

func (router *Router) ListApplications(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	applications, err := router.manager.ListApplications()

	if err != nil {
		return err
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

	return nil
}

func (router *Router) ListRoles(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	lr := request.ListRoles{}
	// get paramter from url query
	if err := getQuery(r, &lr); err != nil {
		return err
	}

	er, err := router.manager.ListRoles(lr.AppID)

	if err != nil {
		return err
	}

	if err := responseOK(er.Format(), w); err != nil {
		router.logger.Error(err.Error())
	}

	return nil
}

func (router *Router) ListScopes(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	ls := request.ListScopes{}
	// get paramter from url query
	if err := getQuery(r, &ls); err != nil {
		return err
	}

	es, err := router.manager.ListScopes(ls.AppID)
	if err != nil {
		return err
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

	return nil
}

func (router *Router) SetUserRole(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	_, perr := getAccessPayload(r)
	if perr != nil {
		return perr
	}

	ur := request.UserRole{}
	if err := getBody(r, &ur); err != nil {
		return perr
	}

	err := router.manager.SetUserRole(ur)

	if err != nil {
		return perr
	}

	if err := responseOK(nil, w); err != nil {
		router.logger.Error(err.Error())
	}

	return nil
}

func (router *Router) FindUser(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	_, perr := getAccessPayload(r)
	if perr != nil {
		return perr
	}
	fu := &request.FindUser{}
	if err := getQuery(r, fu); err != nil {
		return perr
	}

	founds, err := router.manager.FindUser(fu)
	if err != nil {
		return perr
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

	return nil
}

func (router *Router) UsersCount(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	_, perr := getAccessPayload(r)
	if perr != nil {
		return perr
	}
	res, perr := router.manager.UsersCount()
	if perr != nil {
		return perr
	}
	if err := responseOK(res, w); err != nil {
		router.logger.Error(err.Error())
	}

	return nil
}
