package manage

import (
	"database/sql"

	"github.com/volatiletech/sqlboiler/boil"

	"github.com/leeif/pluto/config"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/models"
	"github.com/leeif/pluto/utils/general"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func getUserRole(userID uint, appID string, db boil.Executor) (*models.RbacRole, *perror.PlutoError) {
	app, err := models.Applications(qm.Where("name = ?", appID)).One(db)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err != nil && err == sql.ErrNoRows {
		return nil, perror.ApplicationNotExist
	}

	userAppRole, err := models.RbacUserApplicationRoles(qm.Where("user_id = ? and app_id = ?", userID, app.ID)).One(db)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	var role *models.RbacRole

	if userAppRole != nil {
		role, err = models.RbacRoles(qm.Where("id = ?", userAppRole.RoleID)).One(db)
		if err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	// forbidd pluto admin application with default role assign
	if !app.DefaultRole.IsZero() && app.Name != general.PlutoAdminApplication {
		role, err = models.RbacRoles(qm.Where("id = ?", app.DefaultRole.Uint)).One(db)
		if err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}
	return role, nil
}

func getRoleScopes(role *models.RbacRole, db boil.Executor) (*models.RbacScope, *perror.PlutoError) {

	if role.DefaultScope.IsZero() {
		return nil, perror.ScopeNotExist
	}

	scope, err := models.RbacScopes(qm.Where("id = ?", role.DefaultScope.Uint)).One(db)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	return scope, nil
}

func getUserScopes(userID uint, appID string, db boil.Executor) (*models.RbacScope, *perror.PlutoError) {
	role, err := getUserRole(userID, appID, db)
	if err != nil {
		return nil, err
	}

	if role == nil {
		return nil, nil
	}

	scope, err := getRoleScopes(role, db)

	return scope, nil
}

type Manager struct {
	logger *log.PlutoLog
	config *config.Config
	db     *sql.DB
}

func NewManager(db *sql.DB, config *config.Config, logger *log.PlutoLog) *Manager {
	return &Manager{
		logger: logger.With("compoment", "manager"),
		config: config,
		db:     db,
	}
}
