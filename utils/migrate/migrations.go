package migrate

import "database/sql"

var migrations = []Migrations{
	{
		name:     "create_users_table",
		function: changeUsersTable,
	},
	{
		name:     "create_refresh_tokens_table",
		function: changeRefreshTokensTable,
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
		name:     "create_history_operations_table",
		function: createHistoryOperationsTable,
	},
	{
		name:     "drop_history_operations_table",
		function: dropHistoryOperationsTable,
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
	// {
	// 	name:     "add_identifier_column_for_applications_table",
	// 	function: addIdentifierColumnForApplicationTable,
	// },
	// {
	// 	name:     "add_identifier_index_for_applications_table",
	// 	function: addIdentifierIndexForApplicationTable,
	// },
}

func changeUsersTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `users` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`mail` varchar(255) DEFAULT NULL," +
		"`name` varchar(60) NOT NULL," +
		"`gender` varchar(10) DEFAULT NULL," +
		"`password` varchar(255) DEFAULT NULL," +
		"`birthday` timestamp NULL DEFAULT NULL," +
		"`avatar` varchar(255) DEFAULT NULL," +
		"`verified` tinyint(1) DEFAULT NULL," +
		"`login_type` varchar(10) NOT NULL," +
		"`identify_token` varchar(255) NOT NULL," +
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

func changeRefreshTokensTable(db *sql.DB, name string) error {
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
		"KEY `user_id_refresh_token` (`user_id`,`refresh_token`)" +
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
		"`app_id` varchar(255) NOT NULL," +
		"PRIMARY KEY (`id`)," +
		"KEY `idx_device_apps_deleted_at` (`deleted_at`)," +
		"KEY `device_id_app_id` (`device_id`,`app_id`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createHistoryOperationsTable(db *sql.DB, name string) error {
	sql := "CREATE TABLE IF NOT EXISTS `history_operations` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT," +
		"`created_at` timestamp NULL DEFAULT NULL," +
		"`updated_at` timestamp NULL DEFAULT NULL," +
		"`deleted_at` timestamp NULL DEFAULT NULL," +
		"`user_id` int(10) unsigned NOT NULL," +
		"`type` varchar(20) NOT NULL," +
		"PRIMARY KEY (`id`)," +
		"KEY `idx_history_operations_deleted_at` (`deleted_at`)," +
		"KEY `user_id` (`user_id`)" +
		")"
	_, err := db.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func dropHistoryOperationsTable(db *sql.DB, name string) error {
	sql := "DROP TABLE IF EXISTS `history_operations`"
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
		"KEY `app_id` (`app_id`)," +
		"UNIQUE KEY `name` (`name`)," +
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
		"UNIQUE KEY `name` (`name`)," +
		"KEY `app_id` (`app_id`)" +
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

// func addIdentifierColumnForApplicationTable(db *sql.DB, name string) error {
// 	sql := "ALTER TABLE `applications` ADD `identifier` varchar(255) NOT NULL;"
// 	_, err := db.Exec(sql)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

// func addIdentifierIndexForApplicationTable(db *sql.DB, name string) error {
// 	sql := "CREATE INDEX `identifier` ON `applications` (`identifier`)"
// 	_, err := db.Exec(sql)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
