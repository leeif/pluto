package request

type FindUser struct {
	Keyword string `json:"keyword"`
}

func (fu *FindUser) Validation() bool {
	if fu.Keyword == "" {
		return false
	}
	return true
}

type CreateRole struct {
	AppID uint   `json:"app_id"`
	Name  string `json:"name"`
}

func (cr *CreateRole) Validation() bool {
	if cr.AppID == 0 || cr.Name == "" {
		return false
	}
	return true
}

type CreateScope struct {
	AppID uint   `json:"app_id"`
	Name  string `json:"name"`
}

func (cs *CreateScope) Validation() bool {
	if cs.AppID == 0 || cs.Name == "" {
		return false
	}
	return true
}

type CreateApplication struct {
	Name string `json:"name"`
}

func (ca *CreateApplication) Validation() bool {
	if ca.Name == "" {
		return false
	}
	return true
}

type RoleScope struct {
	RoleID  uint `json:"role_id"`
	ScopeID uint `json:"scope_id"`
}

func (rs *RoleScope) Validation() bool {
	if rs.RoleID == 0 {
		return false
	}
	if rs.ScopeID == 0 {
		return false
	}
	return true
}

type RoleScopeUpdate struct {
	RoleID uint   `json:"role_id"`
	Scopes []uint `json:"scopes"`
}

func (rscu *RoleScopeUpdate) Validation() bool {
	if rscu.RoleID == 0 {
		return false
	}
	if rscu.Scopes == nil || len(rscu.Scopes) == 0 {
		return false
	}
	return true
}

type ApplicationRole struct {
	AppID  uint `json:"app_id"`
	RoleID uint `json:"role_id"`
}

type ListRoles struct {
	AppID uint `json:"app_id" schema:"app_id,required"`
}

type ListScopes struct {
	AppID uint `json:"app_id" schema:"app_id,required"`
}

type UserRole struct {
	UserID uint `json:"user_id" schema:"user_id,required"`
	AppID  uint `json:"app_id" schema:"app_id,required"`
	RoleID uint `json:"role_id" schema:"role_id,required"`
}
