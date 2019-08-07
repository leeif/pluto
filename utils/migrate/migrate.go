package migrate

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/leeif/pluto/models"
)

type Migrations struct {
	name     string
	function func(db *gorm.DB, name string) error
}

func Migrate(db *gorm.DB) error {

	// auto migrate
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Salt{})
	db.AutoMigrate(&models.Device{})
	db.AutoMigrate(&models.RefreshToken{})
	db.AutoMigrate(&models.Service{})

	if err := createMigrationTable(db); err != nil {
		return err
	}

	for _, m := range migrations {
		if migrationNotExists(db, m.name) == false {
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

func createMigrationTable(db *gorm.DB) error {

	if exists := db.HasTable(&models.Migration{}); exists {
		return nil
	}

	// Create migrations table if not exists
	if err := db.CreateTable(&models.Migration{}).Error; err != nil {
		return fmt.Errorf("Error creating migrations table: %s", db.Error)
	}

	return nil
}

func migrationNotExists(db *gorm.DB, name string) bool {
	migration := models.Migration{}
	notFound := db.Where("name = ?", name).First(&migration).RecordNotFound()
	return notFound
}

func saveMigration(db *gorm.DB, name string) error {
	migration := models.Migration{}
	migration.Name = name

	if err := db.Create(&migration).Error; err != nil {
		return fmt.Errorf("Error saving record to migrations table: %s", err)
	}

	return nil
}
