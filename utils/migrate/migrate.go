package migrate

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Migrations struct {
	name     string
	function func(db *sql.DB, name string) error
}

func Migrate(db *sql.DB) error {

	if err := createMigrationTable(db); err != nil {
		return err
	}

	for _, m := range migrations {
		exists, err := migrationExists(db, m.name)

		if err != nil {
			return err
		}

		if exists {
			continue
		}

		if err := m.function(db, m.name); err != nil {
			return err
		}

		if err := saveMigration(db, m.name); err != nil {
			return err
		}
	}

	return nil
}

func createMigrationTable(db *sql.DB) error {

	_, err := db.Exec("CREATE TABLE IF NOT EXISTS `migrations` (" +
		"`id` int(10) unsigned NOT NULL AUTO_INCREMENT, " +
		"`created_at` timestamp NULL DEFAULT NULL, " +
		"`name` varchar(100) NOT NULL, " +
		"PRIMARY KEY (`id`))")
	if err != nil {
		return err
	}

	return nil
}

func migrationExists(db *sql.DB, name string) (bool, error) {
	res, err := db.Query("select * from migrations where name = ?", name)
	if err != nil {
		return false, err
	}
	if res.Next() {
		return true, nil
	}
	return false, nil
}

func saveMigration(db *sql.DB, name string) error {
	log.Println("Start " + name)
	createAt := time.Now()
	res, err := db.Exec("insert into migrations (created_at, name) values (?, ?)", createAt, name)

	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected < 1 {
		return fmt.Errorf("save migration %s failed", name)
	}

	return nil
}
