package manage

import (
	"database/sql"
	"strings"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/wxnacy/wgo/arrays"

	"github.com/leeif/pluto/config"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/log"
	plog "github.com/leeif/pluto/log"
	"github.com/leeif/pluto/models"
	"github.com/leeif/pluto/utils/general"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func getUserRole(exec boil.Executor, userID uint, appID string) (*models.RbacRole, *perror.PlutoError) {
	app, err := models.Applications(qm.Where("name = ?", appID)).One(exec)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err != nil && err == sql.ErrNoRows {
		return nil, perror.ApplicationNotExist
	}

	userAppRole, err := models.RbacUserApplicationRoles(qm.Where("user_id = ? and app_id = ?", userID, app.ID)).One(exec)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	var role *models.RbacRole

	if userAppRole != nil {
		role, err = models.RbacRoles(qm.Where("id = ?", userAppRole.RoleID)).One(exec)
		if err != nil && err != sql.ErrNoRows {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	if role != nil {
		return role, nil
	}

	// forbidd pluto admin application with default role assign
	if !app.DefaultRole.IsZero() && app.Name != general.PlutoAdminApplication {
		role, err = models.RbacRoles(qm.Where("id = ?", app.DefaultRole.Uint)).One(exec)
		if err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	return role, nil
}

func getRoleDefaultScopes(exec boil.Executor, role *models.RbacRole) (*models.RbacScope, *perror.PlutoError) {

	if role.DefaultScope.IsZero() {
		return nil, nil
	}

	scope, err := models.RbacScopes(qm.Where("id = ?", role.DefaultScope.Uint)).One(exec)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err == sql.ErrNoRows {
		return nil, nil
	}

	return scope, nil
}

func getRoleAllScopes(exec boil.Executor, role *models.RbacRole) (models.RbacScopeSlice, *perror.PlutoError) {

	if role == nil || role.DefaultScope.IsZero() {
		return nil, nil
	}

	roleScopes, err := models.RbacRoleScopes(qm.Where("role_id = ?", role.ID)).All(exec)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	roleScopesIn := make([]interface{}, 0)

	for _, roleScope := range roleScopes {
		roleScopesIn = append(roleScopesIn, roleScope.ID)
	}

	scopes, err := models.RbacScopes(qm.WhereIn("id in ?", roleScopesIn...)).All(exec)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	return scopes, nil
}

func getUserDefaultScopes(exec boil.Executor, userID uint, appID string) ([]string, *perror.PlutoError) {
	res := make([]string, 0)
	role, err := getUserRole(exec, userID, appID)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return res, nil
	}

	scope, err := getRoleDefaultScopes(exec, role)

	if err != nil {
		return nil, err
	}

	if scope == nil {
		return res, nil
	}

	res = append(res, scope.Name)

	return res, nil
}

func getUserAllScopes(exec boil.Executor, userID uint, appID string) ([]string, *perror.PlutoError) {
	res := make([]string, 0)
	role, err := getUserRole(exec, userID, appID)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return res, nil
	}

	scopes, err := getRoleAllScopes(exec, role)

	if err != nil {
		return nil, err
	}

	for _, scope := range scopes {
		res = append(res, scope.Name)
	}

	return res, nil
}

func getValidScopes(exec boil.Executor, requestScopes string, userID uint, appID string) (string, *perror.PlutoError) {

	scopesSlice := strings.Split(requestScopes, ",")

	if len(scopesSlice) == 0 {
		return "", nil
	}

	role, perr := getUserRole(exec, userID, appID)
	if perr != nil {
		return "", perr
	}

	scopes, perr := getRoleAllScopes(exec, role)

	if perr != nil {
		return "", perr
	}

	res := make([]string, 0)
	defaultScope := ""
	for _, scope := range scopes {
		if arrays.Contains(scopesSlice, scope.Name) != -1 {
			res = append(res, scope.Name)
		}
		if !role.DefaultScope.IsZero() && scope.ID == role.DefaultScope.Uint {
			defaultScope = scope.Name
		}
	}

	if len(res) == 0 && defaultScope != "" {
		res = append(res, defaultScope)
	}

	return strings.Join(res, ","), nil
}

type Manager struct {
	logger *log.PlutoLog
	config *config.Config
	db     *sql.DB
}

func NewManager(db *sql.DB, config *config.Config, logger *plog.PlutoLog) (*Manager, error) {
	manager := &Manager{
		config: config,
		db:     db,
	}

	if logger == nil {
		var err error
		logger, err = plog.NewLogger(config)
		if err != nil {
			return nil, err
		}
	}

	manager.logger = logger

	return manager, nil
}
