package manage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/volatiletech/sqlboiler/boil"
	"github.com/wxnacy/wgo/arrays"

	"github.com/MuShare/pluto/config"
	perror "github.com/MuShare/pluto/datatype/pluto_error"
	plog "github.com/MuShare/pluto/log"
	"github.com/MuShare/pluto/modelexts"
	"github.com/MuShare/pluto/models"
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

	if !app.DefaultRole.IsZero() {
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
	logger *plog.PlutoLog
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

func getAppAppleLogin(m *Manager, appID string) (*modelexts.AppleLogin, *perror.PlutoError) {
	app, err := models.Applications(qm.Where("name = ?", appID)).One(m.db)

	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	if err == sql.ErrNoRows {
		return nil, perror.ApplicationNotExist
	}

	result := &modelexts.AppleLogin{}
	err = app.AppleLogin.Unmarshal(result)

	if err != nil {
		return nil, perror.ServerError.Wrapper(fmt.Errorf("failed to unmarshal appleLogin: %s", appID))
	}

	return result, nil
}

func getAppWechatLogin(m *Manager, appID string) (*modelexts.WechatLogin, *perror.PlutoError) {
	app, err := models.Applications(qm.Where("name = ?", appID)).One(m.db)

	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	if err == sql.ErrNoRows {
		return nil, perror.ApplicationNotExist
	}

	result := &modelexts.WechatLogin{}
	err = app.WechatLogin.Unmarshal(result)

	if err != nil {
		return nil, perror.ServerError.Wrapper(fmt.Errorf("failed to unmarshal wechatLogin: %s", appID))
	}

	return result, nil
}
