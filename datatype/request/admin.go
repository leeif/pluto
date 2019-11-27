package request

type FindUser struct {
	Keyword string `json:"keyword"`
}

type CreateRole struct {
	AppID uint   `json:"app_id"`
	Name  string `json:"name"`
}

type CreateScope struct {
	AppID uint   `json:"app_id"`
	Name  string `json:"name"`
}

type CreateApplication struct {
	Name string `json:"name"`
}

type RoleScope struct {
	RoleID  uint `json:"role_id"`
	ScopeID uint `json:"scope_id"`
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
