package v1

import (
	"net/http"

	perror "github.com/MuShare/pluto/datatype/pluto_error"

	"github.com/MuShare/pluto/datatype/request"
	routeUtils "github.com/MuShare/pluto/utils/route"
)

func (router *Router) CreateRole(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	cr := request.CreateRole{}
	if perr := routeUtils.GetRequestData(r, &cr); perr != nil {
		return perr
	}
	if _, perr := router.manager.CreateRole(cr); perr != nil {
		return perr
	}
	if err := routeUtils.ResponseOK(nil, w); err != nil {
	}
	return nil
}

func (router *Router) CreateScope(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	cs := request.CreateScope{}
	if perr := routeUtils.GetRequestData(r, &cs); perr != nil {
		return perr
	}
	if _, perr := router.manager.CreateScope(cs); perr != nil {
		return perr
	}
	if err := routeUtils.ResponseOK(nil, w); err != nil {
		return err
	}

	return nil
}

func (router *Router) RoleScopeUpdate(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	rsu := request.RoleScopeUpdate{}
	if perr := routeUtils.GetRequestData(r, &rsu); perr != nil {
		return perr
	}
	if perr := router.manager.RoleScopeUpdate(rsu); perr != nil {
		return perr
	}
	if err := routeUtils.ResponseOK(nil, w); err != nil {
		return err
	}

	return nil
}

func (router *Router) CreateApplication(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	ca := request.CreateApplication{}
	if perr := routeUtils.GetRequestData(r, &ca); perr != nil {
		return perr
	}
	application, perr := router.manager.CreateApplication(ca)
	if perr != nil {
		return perr
	}
	res := make(map[string]interface{})
	res["id"] = application.ID
	res["name"] = application.Name
	res["webhook"] = application.Webhook
	res["defatul_role"] = application.DefaultRole
	if err := routeUtils.ResponseOK(res, w); err != nil {
		return err
	}

	return nil
}

func (router *Router) ApplicationDefaultRole(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	ar := request.ApplicationRole{}
	if perr := routeUtils.GetRequestData(r, &ar); perr != nil {
		return perr
	}
	if perr := router.manager.ApplicationDefaultRole(ar); perr != nil {
		return perr
	}
	if err := routeUtils.ResponseOK(nil, w); err != nil {
		return err
	}

	return nil
}

func (router *Router) RoleDefaultScope(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	rs := request.RoleScope{}
	if perr := routeUtils.GetRequestData(r, &rs); perr != nil {
		return perr
	}
	if perr := router.manager.RoleDefaultScope(rs); perr != nil {
		return perr
	}
	if err := routeUtils.ResponseOK(nil, w); err != nil {
		return err
	}

	return nil
}

func (router *Router) ListApplications(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	applications, perr := router.manager.ListApplications()

	if perr != nil {
		return perr
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

	if err := routeUtils.ResponseOK(res, w); err != nil {
		return err
	}

	return nil
}

func (router *Router) ListRoles(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	lr := request.ListRoles{}
	// get paramter from url query
	if err := routeUtils.GetRequestData(r, &lr); err != nil {
		return err
	}

	er, err := router.manager.ListRoles(lr.AppID)

	if err != nil {
		return err
	}

	if err := routeUtils.ResponseOK(er.Format(), w); err != nil {
		return err
	}

	return nil
}

func (router *Router) ListScopes(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	ls := request.ListScopes{}
	// get paramter from url query
	if err := routeUtils.GetRequestData(r, &ls); err != nil {
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

	if err := routeUtils.ResponseOK(res, w); err != nil {
		return err
	}

	return nil
}

func (router *Router) SetUserRole(w http.ResponseWriter, r *http.Request) *perror.PlutoError {

	ur := request.UserRole{}
	if perr := routeUtils.GetRequestData(r, &ur); perr != nil {
		return perr
	}

	if perr := router.manager.SetUserRole(ur); perr != nil {
		return perr
	}

	if err := routeUtils.ResponseOK(nil, w); err != nil {
		return err
	}

	return nil
}

func (router *Router) FindUser(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	fu := &request.FindUser{}
	if perr := routeUtils.GetRequestData(r, fu); perr != nil {
		return perr
	}

	founds, perr := router.manager.FindUser(fu)
	if perr != nil {
		return perr
	}
	users := make([]map[string]interface{}, 0)
	for _, found := range founds {
		users = append(users, found.Format())
	}
	res := make(map[string]interface{})
	res["users"] = users
	if err := routeUtils.ResponseOK(res, w); err != nil {
		return err
	}

	return nil
}

func (router *Router) UserSummary(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	res, perr := router.manager.UserSummary()
	if perr != nil {
		return perr
	}
	if err := routeUtils.ResponseOK(res, w); err != nil {
		return err
	}

	return nil
}
