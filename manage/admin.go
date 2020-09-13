package manage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/volatiletech/null"

	"github.com/MuShare/pluto/modelexts"

	perror "github.com/MuShare/pluto/datatype/pluto_error"
	"github.com/MuShare/pluto/datatype/request"
	"github.com/MuShare/pluto/models"
	"github.com/volatiletech/sqlboiler/boil"
	"github.com/volatiletech/sqlboiler/queries/qm"
)

func (m *Manager) CreateRole(cr request.CreateRole) (*models.RbacRole, *perror.PlutoError) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	if _, err := models.Applications(qm.Where("id = ?", cr.AppID)).One(tx); err != nil {
		if err == sql.ErrNoRows {
			return nil, perror.ApplicationNotExist
		}
		return nil, perror.ServerError.Wrapper(err)
	}

	role, err := models.RbacRoles(qm.Where("app_id = ? and name = ?", cr.AppID, cr.Name)).One(tx)
	if err == nil {
		return role, perror.RoleExists
	} else if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err != nil && err == sql.ErrNoRows {
		role = &models.RbacRole{}
		role.Name = cr.Name
		role.AppID = cr.AppID
		if err := role.Insert(tx, boil.Infer()); err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	tx.Commit()
	return role, nil
}

func (m *Manager) CreateScope(cs request.CreateScope) (*models.RbacScope, *perror.PlutoError) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	if _, err := models.Applications(qm.Where("id = ?", cs.AppID)).One(tx); err != nil {
		if err == sql.ErrNoRows {
			return nil, perror.ApplicationNotExist
		}
		return nil, perror.ServerError.Wrapper(err)
	}

	scope, err := models.RbacScopes(qm.Where("app_id = ? and name = ?", cs.AppID, cs.Name)).One(tx)
	if err == nil {
		return scope, perror.ScopeExists
	} else if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err != nil && err == sql.ErrNoRows {
		scope = &models.RbacScope{}
		scope.Name = cs.Name
		scope.AppID = cs.AppID
		if err := scope.Insert(tx, boil.Infer()); err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	tx.Commit()
	return scope, nil
}

func (m *Manager) CreateApplication(ca request.CreateApplication) (*models.Application, *perror.PlutoError) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	app, err := models.Applications(qm.Where("name = ?", ca.Name)).One(tx)
	if err == nil {
		return app, perror.ApplicationExists
	} else if err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else {
		app = &models.Application{}
		app.Name = ca.Name
		if err := app.Insert(tx, boil.Infer()); err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	tx.Commit()
	return app, nil
}

func (m *Manager) UpdateApplicationI18nNames(uai request.UpdateApplicationI18Name) (*perror.PlutoError) {
	tx, err := m.db.Begin()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	app, err := models.Applications(qm.Where("id = ?", uai.AppID)).One(tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return perror.ApplicationNotExist
		} else {
			return perror.ServerError.Wrapper(err)
		}
	} else {
		i18nNames, err := json.Marshal(uai.I18Names)
		if err != nil {
			return perror.ServerError.Wrapper(err)
		}
		app.I18nApplicationName = null.JSONFrom(i18nNames)
		if _, err := app.Update(tx, boil.Infer()); err != nil {
			return perror.ServerError.Wrapper(err)
		}
		tx.Commit()
		return nil
	}
}

func (m *Manager) ApplicationI18nName(appName string, language string) (string, *perror.PlutoError) {
	tx, err := m.db.Begin()
	if err != nil {
		return "", perror.ServerError.Wrapper(err)
	}
	defer func() {
		tx.Rollback()
	}()
	app, err := models.Applications(qm.Where("name = ?", appName)).One(tx)
	if app == nil {
		m.logger.Warn(fmt.Sprintf("No i18n name for %s, fallback to pluto", appName))
		return "Pluto", nil
	} else {
		var i18nNames []modelexts.ApplicationI18nName
		err := app.I18nApplicationName.Unmarshal(&i18nNames)
		if err != nil {
			m.logger.Warn(fmt.Sprintf("Failed to unmarshal i18nNames: %s", appName))
			return "pluto", nil
		}

		for _, i18nName := range i18nNames {
			if i18nName.Language == language {
				return i18nName.Name, nil
			}
		}
		return "pluto", nil
	}
}

func (m *Manager) ApplicationI18nNameList(appId uint) (*modelexts.ApplicationI18nNameInfo, *perror.PlutoError) {
	tx, err := m.db.Begin()
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	app, err := models.Applications(qm.Where("id = ?", appId)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	} else if err != nil && err == sql.ErrNoRows {
		return nil, perror.ApplicationNotExist
	}
	var names []modelexts.ApplicationI18nName
	if err := app.I18nApplicationName.Unmarshal(&names); err != nil {
		m.logger.Error(fmt.Sprintf("Failed to unmarshal i18nName, app id: %d", appId))
		return nil, perror.ServerError.Wrapper(err)
	}
	res := modelexts.ApplicationI18nNameInfo{
		AppId:     app.ID,
		AppName:   app.Name,
		I18nNames: &names,
	}
	return &res, nil
}

func (m *Manager) RoleScopeUpdate(rsu request.RoleScopeUpdate) *perror.PlutoError {
	tx, err := m.db.Begin()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	defer func() {
		tx.Rollback()
	}()

	role, err := models.RbacRoles(qm.Where("id = ?", rsu.RoleID)).One(tx)
	if err != nil && err != sql.ErrNoRows {
		return perror.ServerError.Wrapper(err)
	} else if err != nil && err == sql.ErrNoRows {
		return perror.RoleNotExist
	}

	in := make([]interface{}, 0)
	for _, scope := range rsu.Scopes {
		in = append(in, scope)
	}

	scopes, err := models.RbacScopes(qm.WhereIn("id in ?", in...)).All(tx)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	scopeMap := make(map[uint]*models.RbacScope)
	for _, scope := range scopes {
		scopeMap[scope.ID] = scope
	}

	roleScopes, err := models.RbacRoleScopes(qm.Where("role_id = ?", role.ID)).All(m.db)
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	for _, roleScope := range roleScopes {
		if _, ok := scopeMap[roleScope.ScopeID]; ok {
			scopeMap[roleScope.ScopeID] = nil
			continue
		}

		if _, err := roleScope.Delete(tx); err != nil {
			return perror.ServerError.Wrapper(err)
		}
		scopeMap[roleScope.ScopeID] = nil
	}

	for _, v := range scopeMap {
		if v == nil || v.AppID != role.AppID {
			continue
		}
		roleScope := models.RbacRoleScope{}
		roleScope.RoleID = role.ID
		roleScope.ScopeID = v.ID
		if err := roleScope.Insert(tx, boil.Infer()); err != nil {
			return perror.ServerError.Wrapper(err)
		}
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

func (m *Manager) ListRoles(appID uint) (*modelexts.Roles, *perror.PlutoError) {
	application, err := models.Applications(qm.Where("id = ?", appID)).One(m.db)

	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	if err == sql.ErrNoRows {
		return nil, perror.ApplicationNotExist
	}

	roles, err := models.RbacRoles(qm.Where("app_id = ?", appID)).All(m.db)

	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	er := &modelexts.Roles{}
	er.Application = application
	er.Roles = make([]modelexts.Role, 0)

	roleIDs := make([]interface{}, 0)
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	var roleScopes models.RbacRoleScopeSlice

	if len(roleIDs) > 0 {
		roleScopes, err = models.RbacRoleScopes(qm.WhereIn("role_id in ?", roleIDs...)).All(m.db)
		if err != nil && err != sql.ErrNoRows {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	scopeIDs := make([]interface{}, 0)
	for _, k := range roleScopes {
		scopeIDs = append(scopeIDs, k.ScopeID)
	}

	var scopes models.RbacScopeSlice

	if len(scopeIDs) > 0 {
		scopes, err = models.RbacScopes(qm.WhereIn("id in ?", scopeIDs...)).All(m.db)
		if err != nil && err != sql.ErrNoRows {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	scopeMap := make(map[uint]*models.RbacScope)
	for _, scope := range scopes {
		scopeMap[scope.ID] = scope
	}

	roleScopeMap := make(map[uint][]*models.RbacScope)

	for _, rs := range roleScopes {
		if _, ok := roleScopeMap[rs.RoleID]; !ok {
			roleScopeMap[rs.RoleID] = make([]*models.RbacScope, 0)
		}
		roleScopeMap[rs.RoleID] = append(roleScopeMap[rs.RoleID], scopeMap[rs.ScopeID])
	}

	for _, role := range roles {
		r := modelexts.Role{}
		r.RbacRole = role
		if scopes, ok := roleScopeMap[r.ID]; ok {
			r.Scopes = scopes
		}
		er.Roles = append(er.Roles, r)
	}

	return er, nil
}

func (m *Manager) ListScopes(appID uint) (*modelexts.Scopes, *perror.PlutoError) {
	application, err := models.Applications(qm.Where("id = ?", appID)).One(m.db)

	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	if err == sql.ErrNoRows {
		return nil, perror.ApplicationNotExist
	}

	scopes, err := models.RbacScopes(qm.Where("app_id = ?", appID)).All(m.db)

	if err != nil && err != sql.ErrNoRows {
		return nil, perror.ServerError.Wrapper(err)
	}

	es := &modelexts.Scopes{}
	es.Application = application
	es.Scopes = scopes

	return es, nil
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
		tx.Commit()
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

func (m *Manager) UserSummary() (map[string]int, *perror.PlutoError) {
	res := make(map[string]int)
	users, err := models.Users().All(m.db)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}
	res["total"] = len(users)
	return res, nil
}

func (m *Manager) FindUser(fu *request.FindUser) ([]*modelexts.FindUser, *perror.PlutoError) {

	users, err := models.Users(qm.Where("name = ?", fu.Name)).All(m.db)

	if err != nil && len(users) == 0 {
		return nil, perror.UserNotExist
	} else if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	res := make([]*modelexts.FindUser, 0)

	for _, user := range users {
		found, err := m.findUser(user)
		if err != nil {
			return nil, err
		}
		res = append(res, found)
	}

	return res, nil
}

func (m *Manager) findUser(user *models.User) (*modelexts.FindUser, *perror.PlutoError) {
	applications, err := models.Applications().All(m.db)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	appMap := make(map[uint]*models.Application)

	for _, application := range applications {
		appMap[application.ID] = application
	}

	userAppRoles, err := models.RbacUserApplicationRoles(qm.Where("user_id = ?", user.ID)).All(m.db)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	ra := make([]interface{}, 0)
	for _, userAppRole := range userAppRoles {
		ra = append(ra, userAppRole.RoleID)
	}

	var roles models.RbacRoleSlice
	if len(ra) > 0 {
		roles, err = models.RbacRoles(qm.WhereIn("id in ?", ra...)).All(m.db)
		if err != nil {
			return nil, perror.ServerError.Wrapper(err)
		}
	}

	roleMap := make(map[uint]*models.RbacRole)
	for _, role := range roles {
		roleMap[role.ID] = role
	}

	userAppRoleMap := make(map[uint]*modelexts.UserApplicationRole)
	for _, userAppRole := range userAppRoles {
		if _, ok := userAppRoleMap[userAppRole.AppID]; !ok {
			userAppRoleMap[userAppRole.AppID] = &modelexts.UserApplicationRole{}
			userAppRoleMap[userAppRole.AppID].Application = appMap[userAppRole.AppID]
			userAppRoleMap[userAppRole.AppID].Roles = make([]*models.RbacRole, 0)
		}

		if role, ok := roleMap[userAppRole.RoleID]; ok {
			userAppRoleMap[userAppRole.AppID].Roles = append(userAppRoleMap[userAppRole.AppID].Roles, role)
		}
	}

	extApps := make([]*modelexts.UserApplicationRole, 0)

	for _, application := range applications {
		extApp := &modelexts.UserApplicationRole{}
		if uar, ok := userAppRoleMap[application.ID]; ok {
			extApp = uar
		} else {
			extApp.Application = application
		}
		extApps = append(extApps, extApp)
	}

	bindings, err := models.Bindings(qm.Where("user_id = ?", user.ID)).All(m.db)
	if err != nil {
		return nil, perror.ServerError.Wrapper(err)
	}

	res := &modelexts.FindUser{
		User:         user,
		Bindings:     bindings,
		Applications: extApps,
	}

	return res, nil
}
