package manage

import (
	"database/sql"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (m *Manager) CreateRole(cr request.CreateRole) *perror.PlutoError {
	tx, err := m.db.Begin()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	role := &models.RbacRole{}
	role.Name = cr.Name
	role.AppID = cr.AppID

	if err := role.Insert(tx, boil.Infer()); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	tx.Commit()
	return nil
}

func (m *Manager) CreateScope(cs request.CreateScope) *perror.PlutoError {
	tx, err := m.db.Begin()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	scope, err := models.RbacScopes(qm.Where("name = ?", APPLELOGIN, cs.Name)).One(tx)

	if err != nil && err != sql.ErrNoRows {
		return perror.ServerError.Wrapper(err)
	}

	if err == nil {
		return perror.ScopeExists
	}

	scope = &models.RbacScope{}
	scope.Name = cs.Name
	scope.AppID = cs.AppID

	if err := scope.Insert(tx, boil.Infer()); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	tx.Commit()
	return nil
}

func (m *Manager) CreateApplication(ca request.CreateApplication) (*models.Application, *perror.PlutoError) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	app := &models.Application{}
	app.Name = ca.Name
	app.Identifier = ca.Name + "." + randomToken(6)

	if err := app.Insert(tx, boil.Infer()); err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	tx.Commit()

	return app, nil
}

func (m *Manager) AttachScope(rs request.RoleScope) *perror.PlutoError {
	tx, err := m.db.Begin()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	roleScope, err := models.RbacRoleScopes(qm.Where("role_id = ? and scope_id = ?", rs.RoleID, rs.ScopeID)).One(tx)

	if err != nil && err != sql.ErrNoRows {
		return perror.ServerError.Wrapper(err)
	}

	if err == nil {
		return perror.ScopeAttached
	}

	roleScope = &models.RbacRoleScope{}
	roleScope.RoleID = rs.RoleID
	roleScope.ScopeID = rs.ScopeID

	if err := roleScope.Insert(tx, boil.Infer()); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	tx.Commit()

	return nil
}

func (m *Manager) DetachScope(rs request.RoleScope) *perror.PlutoError {
	tx, err := m.db.Begin()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	roleScope := &models.RbacRoleScope{}
	roleScope.RoleID = rs.RoleID
	roleScope.ScopeID = rs.ScopeID

	if _, err := roleScope.Delete(tx); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	tx.Commit()

	return nil
}

func (m *Manager) ApplicationDefaultRole(ar request.ApplicationRole) *perror.PlutoError {
	tx, err := m.db.Begin()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	application, err := models.Applications(qm.Where("id = ?", ar.AppID)).One(tx)

	if err != nil && err == sql.ErrNoRows {
		return perror.ApplicationNotExist
	} else if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	application.DefaultRole.SetValid(ar.RoleID)

	if _, err := application.Update(tx, boil.Infer()); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	tx.Commit()
	return nil
}

func (m *Manager) RoleDefaultScope(rs request.RoleScope) *perror.PlutoError {
	tx, err := m.db.Begin()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	role, err := models.RbacRoles(qm.Where("id = ?", rs.RoleID)).One(tx)

	if err != nil && err == sql.ErrNoRows {
		return perror.RoleNotExist
	} else if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	scope, err := models.RbacScopes(qm.Where("id = ?", rs.ScopeID)).One(tx)

	if err != nil && err == sql.ErrNoRows {
		return perror.ScopeNotExist
	} else if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	if scope.AppID != role.AppID {
		return perror.ScopeNotExist
	}

	role.DefaultScope.SetValid(rs.ScopeID)
	if _, err := role.Update(tx, boil.Infer()); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	tx.Commit()
	return nil
}

func (m *Manager) ListApplications() (models.ApplicationSlice, *perror.PlutoError) {

	applications, err := models.Applications().All(m.db)

	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	return applications, nil
}

func (m *Manager) ListRoles(appID uint) (models.RbacRoleSlice, *perror.PlutoError) {
	roles, err := models.RbacRoles(qm.Where("app_id = ?", appID)).All(m.db)

	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	return roles, nil
}

func (m *Manager) ListScopes(appID uint) (models.RbacScopeSlice, *perror.PlutoError) {
	scopes, err := models.RbacScopes(qm.Where("app_id = ?", appID)).All(m.db)

	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	return scopes, nil
}

func (m *Manager) SetUserRole(ur request.UserRole) *perror.PlutoError {
	tx, err := m.db.Begin()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	userAppRole, err := models.RbacUserApplicationRoles(qm.Where("user_id = ? and app_id = ?", ur.UserID, ur.AppID)).One(tx)

	if err != nil && err != sql.ErrNoRows {
		return perror.ServerError.Wrapper(err)
	} else if err != nil && err == sql.ErrNoRows {
		userAppRole = &models.RbacUserApplicationRole{}
		userAppRole.UserID = ur.UserID
		userAppRole.AppID = ur.AppID
		userAppRole.RoleID = ur.RoleID
		if err := userAppRole.Insert(tx, boil.Infer()); err != nil {
			return perror.ServerError.Wrapper(err)
		}
		return nil
	}

	role, err := models.RbacRoles(qm.Where("id = ?", ur.RoleID)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return perror.ServerError.Wrapper(err)
	} else if err != nil && err == sql.ErrNoRows {
		return perror.RoleNotExist
	}

	if role.AppID != ur.AppID {
		return perror.RoleNotExist
	}

	userAppRole.RoleID = ur.RoleID
	if _, err := userAppRole.Update(tx, boil.Infer()); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	tx.Commit()
	return nil
}
