package migrate

import "database/sql"

var migrations = []Migrations{
	{
		name:     "create_users_table",
		function: createUsersTable,
	},
	{
		name:     "create_bindings_table",
		function: createBindingsTable,
	},
	{
		name:     "create_refresh_tokens_table",
		function: createRefreshTokensTable,
	},
	{
		name:     "create_salts_table",
		function: createSaltsTable,
	},
	{
		name:     "create_applictaions_table",
		function: createApplicationsTable,
	},
	{
		name:     "create_device_apps_table",
		function: createDeviceAppsTable,
	},
	{
		name:     "create_rbac_user_application_roles_table",
		function: createRBACUserApplicationRoleTable,
	},
	{
		name:     "create_rbac_roles_table",
		function: createRBACRoleTable,
	},
	{
		name:     "create_rbac_role_scopes_table",
		function: createRBACRoleScopeTable,
	},
	{
		name:     "create_rbac_scopes_table",
		function: createRBACScopeTable,
	},
	{
		name:     "add_default_role_column_in_applications_table",
		function: addDefaultRoleColumnInApplicationTable,
	},
	{
		name:     "add_scopes_column_in_refresh_token_table",
		function: addScopesColumnInRefreshTokenTable,
	},
	{
		name:     "create_oauth_clients_table",
		function: createOauthClientsTable,
	},
	{
		name:     "create_oauth_authorization_codes_table",
		function: createOauthAuthorizationCodesTable,
	},
}

func createUsersTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `users` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`name` varchar(60) NOT NULL," +
		"`password` varchar(255) DEFAULT NULL," +
		"`verified` tinyint(1) DEFAULT NULL," +
		"`avatar` varchar(255) DEFAULT NULL," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `name_key` (`name`)," +
		"KEY `users_deleted_at` (`deleted_at`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createBindingsTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `bindings` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`login_type` varchar(10) NOT NULL," +
		"`identify_token` varchar(255) NOT NULL," +
		"`mail` varchar(255) NOT NULL," +
		"`user_id` int(10) unsigned NOT NULL," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `login_type_identify_token` (`login_type`,`identify_token`)," +
		"KEY `users_deleted_at` (`deleted_at`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createRefreshTokensTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `refresh_tokens` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`user_id` int(10) unsigned NOT NULL," +
		"`refresh_token` varchar(255) NOT NULL," +
		"`device_app_id` int(10) unsigned NOT NULL," +
		"PRIMARY KEY (`id`)," +
		"KEY `refresh_tokens_deleted_at` (`deleted_at`)," +
		"KEY `refresh_token` (`refresh_token`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createSaltsTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `salts` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`user_id` int(10) unsigned NOT NULL," +
		"`salt` varchar(255) NOT NULL," +
		"PRIMARY KEY (`id`)," +
		"KEY `user_id` (`user_id`)," +
		"KEY `salts_deleted_at` (`deleted_at`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createApplicationsTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `applications` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`name` varchar(100) NOT NULL," +
		"`webhook` varchar(255) NOT NULL," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `name` (`name`)," +
		"KEY `applictaions_deleted_at` (`deleted_at`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createDeviceAppsTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `device_apps` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`device_id` varchar(255) NOT NULL," +
		"`app_id` int(10) unsigned NOT NULL," +
		"PRIMARY KEY (`id`)," +
		"KEY `idx_device_apps_deleted_at` (`deleted_at`)," +
		"KEY `app_id_device_id` (`app_id`,`device_id`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createRBACUserApplicationRoleTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `rbac_user_application_roles` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`user_id` int(10) unsigned NOT NULL," +
		"`app_id` int(10) unsigned NOT NULL," +
		"`role_id` int(10) unsigned NOT NULL," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `user_id_app_id` (`user_id`, `app_id`)," +
		"KEY `idx_user_roles_deleted_at` (`deleted_at`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createRBACRoleTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `rbac_roles` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`name` varchar(20) NOT NULL," +
		"`app_id` int(10) unsigned NOT NULL," +
		"`default_scope` int(10) unsigned," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `app_id_name` (`app_id`, `name`)," +
		"KEY `idx_roles_deleted_at` (`deleted_at`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createRBACRoleScopeTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `rbac_role_scopes` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`role_id` int(10) unsigned NOT NULL," +
		"`scope_id` int(10) unsigned NOT NULL," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `role_id_scope_id` (`role_id`, `scope_id`)," +
		"KEY `idx_roles_deleted_at` (`deleted_at`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createRBACScopeTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `rbac_scopes` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`name` varchar(20) NOT NULL," +
		"`app_id` int(10) unsigned NOT NULL," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `app_id_name` (`app_id`, `name`)," +
		"KEY `idx_scopes_deleted_at` (`deleted_at`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func addDefaultRoleColumnInApplicationTable(db *sql.DB, name string) error {
	sql := "ALTER TABLE `applications` ADD `default_role` int(10) unsigned;"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func addScopesColumnInRefreshTokenTable(db *sql.DB, name string) error {
	sql := "ALTER TABLE `refresh_tokens` ADD `scopes` varchar(255);"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func removeGenderAndBirthdayColumnInUserTable(db *sql.DB, name string) error {
	sql := "ALTER TABLE `users` DROP COLUMN `gender`, DROP COLUMN `birthday`;"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createOauthClientsTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `oauth_clients` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`key` varchar(255) NOT NULL," +
		"`secret` varchar(60) NOT NULL," +
		"`status` varchar(60) NOT NULL," +
		"`redirect_uri` varchar(200) NOT NULL," +
		"PRIMARY KEY (`id`)," +
		"UNIQUE KEY `client_key` (`key`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createOauthAuthorizationCodesTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `oauth_authorization_codes` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`user_id` int(10) unsigned NOT NULL," +
		"`client_id` int(10) unsigned NOT NULL," +
		"`app_id` int(10) unsigned NOT NULL," +
		"`code` varchar(40) NOT NULL," +
		"`redirect_uri` varchar(200) NOT NULL," +
		"`expire_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP," +
		"`scopes` varchar(200) NOT NULL," +
		"PRIMARY KEY (`id`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
