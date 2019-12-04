package modelexts

import (
	"github.com/leeif/pluto/models"
)

type User struct {
	*models.User
	Roles []string `json:"roles"`
}

type Scopes struct {
	Application *models.Application `json:"application"`
	Scopes      []*models.RbacScope `json:"scopes"`
}

func (scopes Scopes) Format() map[string]interface{} {
	s := make([]map[string]interface{}, 0)

	app := make(map[string]interface{})
	app["id"] = scopes.Application.ID
	app["name"] = scopes.Application.Name
	app["default_role"] = scopes.Application.DefaultRole

	for _, scope := range scopes.Scopes {
		m := make(map[string]interface{})
		m["id"] = scope.ID
		m["name"] = scope.Name
		s = append(s, m)
	}

	res := make(map[string]interface{})
	res["application"] = app
	res["scopes"] = s
	return res
}

type Roles struct {
	Application *models.Application `json:"application"`
	Roles       []Role              `json:"roles"`
}

type Role struct {
	*models.RbacRole
	Scopes []*models.RbacScope `json:"scopes"`
}

func (roles Roles) Format() map[string]interface{} {
	res := make(map[string]interface{})
	app := make(map[string]interface{})
	app["name"] = roles.Application.Name
	app["id"] = roles.Application.ID
	app["default_role"] = roles.Application.DefaultRole

	res["application"] = app

	rs := make([]interface{}, 0)
	for _, role := range roles.Roles {
		r := make(map[string]interface{})
		r["id"] = role.ID
		r["name"] = role.Name
		r["default_scope"] = role.DefaultScope

		scopes := make([]interface{}, 0)
		for _, scope := range role.Scopes {
			s := make(map[string]interface{})
			s["id"] = scope.ID
			s["name"] = scope.Name
			scopes = append(scopes, s)
		}
		r["scopes"] = scopes
		rs = append(rs, r)
	}

	res["roles"] = rs

	return res
}

type UserApplicationRole struct {
	*models.Application
	Roles []*models.RbacRole `json:"roles"`
}

type FindUser struct {
	*models.User
	Applications []*UserApplicationRole `json:"applications"`
}

func (findUser FindUser) Format() map[string]interface{} {
	res := make(map[string]interface{})
	res["id"] = findUser.ID
	res["name"] = findUser.Name
	res["mail"] = findUser.Mail.String
	res["avatar"] = findUser.Avatar.String
	res["login_type"] = findUser.LoginType

	applications := make([]interface{}, 0)
	for _, application := range findUser.Applications {
		a := make(map[string]interface{})
		a["id"] = application.ID
		a["name"] = application.Name
		rs := make([]interface{}, 0)
		for _, role := range application.Roles {
			r := make(map[string]interface{})
			r["id"] = role.Name
			r["name"] = role.ID
			rs = append(rs, r)
		}
		a["roles"] = rs
		applications = append(applications, a)
	}

	res["applications"] = applications

	return res
}
