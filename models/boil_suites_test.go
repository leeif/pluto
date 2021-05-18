// Code generated by SQLBoiler 3.6.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package models

import "testing"

// This test suite runs each operation test in parallel.
// Example, if your database has 3 tables, the suite will run:
// table1, table2 and table3 Delete in parallel
// table1, table2 and table3 Insert in parallel, and so forth.
// It does NOT run each operation group in parallel.
// Separating the tests thusly grants avoidance of Postgres deadlocks.
func TestParent(t *testing.T) {
	t.Run("Applications", testApplications)
	t.Run("Bindings", testBindings)
	t.Run("DeviceApps", testDeviceApps)
	t.Run("Migrations", testMigrations)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodes)
	t.Run("OauthClients", testOauthClients)
	t.Run("RbacRoleScopes", testRbacRoleScopes)
	t.Run("RbacRoles", testRbacRoles)
	t.Run("RbacScopes", testRbacScopes)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRoles)
	t.Run("RefreshTokens", testRefreshTokens)
	t.Run("Salts", testSalts)
	t.Run("Users", testUsers)
}

func TestDelete(t *testing.T) {
	t.Run("Applications", testApplicationsDelete)
	t.Run("Bindings", testBindingsDelete)
	t.Run("DeviceApps", testDeviceAppsDelete)
	t.Run("Migrations", testMigrationsDelete)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesDelete)
	t.Run("OauthClients", testOauthClientsDelete)
	t.Run("RbacRoleScopes", testRbacRoleScopesDelete)
	t.Run("RbacRoles", testRbacRolesDelete)
	t.Run("RbacScopes", testRbacScopesDelete)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesDelete)
	t.Run("RefreshTokens", testRefreshTokensDelete)
	t.Run("Salts", testSaltsDelete)
	t.Run("Users", testUsersDelete)
}

func TestQueryDeleteAll(t *testing.T) {
	t.Run("Applications", testApplicationsQueryDeleteAll)
	t.Run("Bindings", testBindingsQueryDeleteAll)
	t.Run("DeviceApps", testDeviceAppsQueryDeleteAll)
	t.Run("Migrations", testMigrationsQueryDeleteAll)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesQueryDeleteAll)
	t.Run("OauthClients", testOauthClientsQueryDeleteAll)
	t.Run("RbacRoleScopes", testRbacRoleScopesQueryDeleteAll)
	t.Run("RbacRoles", testRbacRolesQueryDeleteAll)
	t.Run("RbacScopes", testRbacScopesQueryDeleteAll)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesQueryDeleteAll)
	t.Run("RefreshTokens", testRefreshTokensQueryDeleteAll)
	t.Run("Salts", testSaltsQueryDeleteAll)
	t.Run("Users", testUsersQueryDeleteAll)
}

func TestSliceDeleteAll(t *testing.T) {
	t.Run("Applications", testApplicationsSliceDeleteAll)
	t.Run("Bindings", testBindingsSliceDeleteAll)
	t.Run("DeviceApps", testDeviceAppsSliceDeleteAll)
	t.Run("Migrations", testMigrationsSliceDeleteAll)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesSliceDeleteAll)
	t.Run("OauthClients", testOauthClientsSliceDeleteAll)
	t.Run("RbacRoleScopes", testRbacRoleScopesSliceDeleteAll)
	t.Run("RbacRoles", testRbacRolesSliceDeleteAll)
	t.Run("RbacScopes", testRbacScopesSliceDeleteAll)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesSliceDeleteAll)
	t.Run("RefreshTokens", testRefreshTokensSliceDeleteAll)
	t.Run("Salts", testSaltsSliceDeleteAll)
	t.Run("Users", testUsersSliceDeleteAll)
}

func TestExists(t *testing.T) {
	t.Run("Applications", testApplicationsExists)
	t.Run("Bindings", testBindingsExists)
	t.Run("DeviceApps", testDeviceAppsExists)
	t.Run("Migrations", testMigrationsExists)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesExists)
	t.Run("OauthClients", testOauthClientsExists)
	t.Run("RbacRoleScopes", testRbacRoleScopesExists)
	t.Run("RbacRoles", testRbacRolesExists)
	t.Run("RbacScopes", testRbacScopesExists)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesExists)
	t.Run("RefreshTokens", testRefreshTokensExists)
	t.Run("Salts", testSaltsExists)
	t.Run("Users", testUsersExists)
}

func TestFind(t *testing.T) {
	t.Run("Applications", testApplicationsFind)
	t.Run("Bindings", testBindingsFind)
	t.Run("DeviceApps", testDeviceAppsFind)
	t.Run("Migrations", testMigrationsFind)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesFind)
	t.Run("OauthClients", testOauthClientsFind)
	t.Run("RbacRoleScopes", testRbacRoleScopesFind)
	t.Run("RbacRoles", testRbacRolesFind)
	t.Run("RbacScopes", testRbacScopesFind)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesFind)
	t.Run("RefreshTokens", testRefreshTokensFind)
	t.Run("Salts", testSaltsFind)
	t.Run("Users", testUsersFind)
}

func TestBind(t *testing.T) {
	t.Run("Applications", testApplicationsBind)
	t.Run("Bindings", testBindingsBind)
	t.Run("DeviceApps", testDeviceAppsBind)
	t.Run("Migrations", testMigrationsBind)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesBind)
	t.Run("OauthClients", testOauthClientsBind)
	t.Run("RbacRoleScopes", testRbacRoleScopesBind)
	t.Run("RbacRoles", testRbacRolesBind)
	t.Run("RbacScopes", testRbacScopesBind)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesBind)
	t.Run("RefreshTokens", testRefreshTokensBind)
	t.Run("Salts", testSaltsBind)
	t.Run("Users", testUsersBind)
}

func TestOne(t *testing.T) {
	t.Run("Applications", testApplicationsOne)
	t.Run("Bindings", testBindingsOne)
	t.Run("DeviceApps", testDeviceAppsOne)
	t.Run("Migrations", testMigrationsOne)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesOne)
	t.Run("OauthClients", testOauthClientsOne)
	t.Run("RbacRoleScopes", testRbacRoleScopesOne)
	t.Run("RbacRoles", testRbacRolesOne)
	t.Run("RbacScopes", testRbacScopesOne)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesOne)
	t.Run("RefreshTokens", testRefreshTokensOne)
	t.Run("Salts", testSaltsOne)
	t.Run("Users", testUsersOne)
}

func TestAll(t *testing.T) {
	t.Run("Applications", testApplicationsAll)
	t.Run("Bindings", testBindingsAll)
	t.Run("DeviceApps", testDeviceAppsAll)
	t.Run("Migrations", testMigrationsAll)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesAll)
	t.Run("OauthClients", testOauthClientsAll)
	t.Run("RbacRoleScopes", testRbacRoleScopesAll)
	t.Run("RbacRoles", testRbacRolesAll)
	t.Run("RbacScopes", testRbacScopesAll)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesAll)
	t.Run("RefreshTokens", testRefreshTokensAll)
	t.Run("Salts", testSaltsAll)
	t.Run("Users", testUsersAll)
}

func TestCount(t *testing.T) {
	t.Run("Applications", testApplicationsCount)
	t.Run("Bindings", testBindingsCount)
	t.Run("DeviceApps", testDeviceAppsCount)
	t.Run("Migrations", testMigrationsCount)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesCount)
	t.Run("OauthClients", testOauthClientsCount)
	t.Run("RbacRoleScopes", testRbacRoleScopesCount)
	t.Run("RbacRoles", testRbacRolesCount)
	t.Run("RbacScopes", testRbacScopesCount)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesCount)
	t.Run("RefreshTokens", testRefreshTokensCount)
	t.Run("Salts", testSaltsCount)
	t.Run("Users", testUsersCount)
}

func TestHooks(t *testing.T) {
	t.Run("Applications", testApplicationsHooks)
	t.Run("Bindings", testBindingsHooks)
	t.Run("DeviceApps", testDeviceAppsHooks)
	t.Run("Migrations", testMigrationsHooks)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesHooks)
	t.Run("OauthClients", testOauthClientsHooks)
	t.Run("RbacRoleScopes", testRbacRoleScopesHooks)
	t.Run("RbacRoles", testRbacRolesHooks)
	t.Run("RbacScopes", testRbacScopesHooks)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesHooks)
	t.Run("RefreshTokens", testRefreshTokensHooks)
	t.Run("Salts", testSaltsHooks)
	t.Run("Users", testUsersHooks)
}

func TestInsert(t *testing.T) {
	t.Run("Applications", testApplicationsInsert)
	t.Run("Applications", testApplicationsInsertWhitelist)
	t.Run("Bindings", testBindingsInsert)
	t.Run("Bindings", testBindingsInsertWhitelist)
	t.Run("DeviceApps", testDeviceAppsInsert)
	t.Run("DeviceApps", testDeviceAppsInsertWhitelist)
	t.Run("Migrations", testMigrationsInsert)
	t.Run("Migrations", testMigrationsInsertWhitelist)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesInsert)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesInsertWhitelist)
	t.Run("OauthClients", testOauthClientsInsert)
	t.Run("OauthClients", testOauthClientsInsertWhitelist)
	t.Run("RbacRoleScopes", testRbacRoleScopesInsert)
	t.Run("RbacRoleScopes", testRbacRoleScopesInsertWhitelist)
	t.Run("RbacRoles", testRbacRolesInsert)
	t.Run("RbacRoles", testRbacRolesInsertWhitelist)
	t.Run("RbacScopes", testRbacScopesInsert)
	t.Run("RbacScopes", testRbacScopesInsertWhitelist)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesInsert)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesInsertWhitelist)
	t.Run("RefreshTokens", testRefreshTokensInsert)
	t.Run("RefreshTokens", testRefreshTokensInsertWhitelist)
	t.Run("Salts", testSaltsInsert)
	t.Run("Salts", testSaltsInsertWhitelist)
	t.Run("Users", testUsersInsert)
	t.Run("Users", testUsersInsertWhitelist)
}

// TestToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestToOne(t *testing.T) {
	t.Run("BindingToApplicationUsingApp", testBindingToOneApplicationUsingApp)
	t.Run("UserToApplicationUsingApp", testUserToOneApplicationUsingApp)
}

// TestOneToOne tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOne(t *testing.T) {}

// TestToMany tests cannot be run in parallel
// or deadlocks can occur.
func TestToMany(t *testing.T) {
	t.Run("ApplicationToAppBindings", testApplicationToManyAppBindings)
	t.Run("ApplicationToAppUsers", testApplicationToManyAppUsers)
}

// TestToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneSet(t *testing.T) {
	t.Run("BindingToApplicationUsingAppBindings", testBindingToOneSetOpApplicationUsingApp)
	t.Run("UserToApplicationUsingAppUsers", testUserToOneSetOpApplicationUsingApp)
}

// TestToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToOneRemove(t *testing.T) {}

// TestOneToOneSet tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneSet(t *testing.T) {}

// TestOneToOneRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestOneToOneRemove(t *testing.T) {}

// TestToManyAdd tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyAdd(t *testing.T) {
	t.Run("ApplicationToAppBindings", testApplicationToManyAddOpAppBindings)
	t.Run("ApplicationToAppUsers", testApplicationToManyAddOpAppUsers)
}

// TestToManySet tests cannot be run in parallel
// or deadlocks can occur.
func TestToManySet(t *testing.T) {}

// TestToManyRemove tests cannot be run in parallel
// or deadlocks can occur.
func TestToManyRemove(t *testing.T) {}

func TestReload(t *testing.T) {
	t.Run("Applications", testApplicationsReload)
	t.Run("Bindings", testBindingsReload)
	t.Run("DeviceApps", testDeviceAppsReload)
	t.Run("Migrations", testMigrationsReload)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesReload)
	t.Run("OauthClients", testOauthClientsReload)
	t.Run("RbacRoleScopes", testRbacRoleScopesReload)
	t.Run("RbacRoles", testRbacRolesReload)
	t.Run("RbacScopes", testRbacScopesReload)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesReload)
	t.Run("RefreshTokens", testRefreshTokensReload)
	t.Run("Salts", testSaltsReload)
	t.Run("Users", testUsersReload)
}

func TestReloadAll(t *testing.T) {
	t.Run("Applications", testApplicationsReloadAll)
	t.Run("Bindings", testBindingsReloadAll)
	t.Run("DeviceApps", testDeviceAppsReloadAll)
	t.Run("Migrations", testMigrationsReloadAll)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesReloadAll)
	t.Run("OauthClients", testOauthClientsReloadAll)
	t.Run("RbacRoleScopes", testRbacRoleScopesReloadAll)
	t.Run("RbacRoles", testRbacRolesReloadAll)
	t.Run("RbacScopes", testRbacScopesReloadAll)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesReloadAll)
	t.Run("RefreshTokens", testRefreshTokensReloadAll)
	t.Run("Salts", testSaltsReloadAll)
	t.Run("Users", testUsersReloadAll)
}

func TestSelect(t *testing.T) {
	t.Run("Applications", testApplicationsSelect)
	t.Run("Bindings", testBindingsSelect)
	t.Run("DeviceApps", testDeviceAppsSelect)
	t.Run("Migrations", testMigrationsSelect)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesSelect)
	t.Run("OauthClients", testOauthClientsSelect)
	t.Run("RbacRoleScopes", testRbacRoleScopesSelect)
	t.Run("RbacRoles", testRbacRolesSelect)
	t.Run("RbacScopes", testRbacScopesSelect)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesSelect)
	t.Run("RefreshTokens", testRefreshTokensSelect)
	t.Run("Salts", testSaltsSelect)
	t.Run("Users", testUsersSelect)
}

func TestUpdate(t *testing.T) {
	t.Run("Applications", testApplicationsUpdate)
	t.Run("Bindings", testBindingsUpdate)
	t.Run("DeviceApps", testDeviceAppsUpdate)
	t.Run("Migrations", testMigrationsUpdate)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesUpdate)
	t.Run("OauthClients", testOauthClientsUpdate)
	t.Run("RbacRoleScopes", testRbacRoleScopesUpdate)
	t.Run("RbacRoles", testRbacRolesUpdate)
	t.Run("RbacScopes", testRbacScopesUpdate)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesUpdate)
	t.Run("RefreshTokens", testRefreshTokensUpdate)
	t.Run("Salts", testSaltsUpdate)
	t.Run("Users", testUsersUpdate)
}

func TestSliceUpdateAll(t *testing.T) {
	t.Run("Applications", testApplicationsSliceUpdateAll)
	t.Run("Bindings", testBindingsSliceUpdateAll)
	t.Run("DeviceApps", testDeviceAppsSliceUpdateAll)
	t.Run("Migrations", testMigrationsSliceUpdateAll)
	t.Run("OauthAuthorizationCodes", testOauthAuthorizationCodesSliceUpdateAll)
	t.Run("OauthClients", testOauthClientsSliceUpdateAll)
	t.Run("RbacRoleScopes", testRbacRoleScopesSliceUpdateAll)
	t.Run("RbacRoles", testRbacRolesSliceUpdateAll)
	t.Run("RbacScopes", testRbacScopesSliceUpdateAll)
	t.Run("RbacUserApplicationRoles", testRbacUserApplicationRolesSliceUpdateAll)
	t.Run("RefreshTokens", testRefreshTokensSliceUpdateAll)
	t.Run("Salts", testSaltsSliceUpdateAll)
	t.Run("Users", testUsersSliceUpdateAll)
}
