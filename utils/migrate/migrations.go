package migrate

import (
	"github.com/jinzhu/gorm"
	"github.com/leeif/pluto/models"
)

var migrations = []Migrations{
	{
		name:     "inital",
		function: migrateInit,
	},
}

func migrateInit(db *gorm.DB, name string) error {
	if err := db.CreateTable(&models.User{}).Error; err != nil {
		return err
	}
	if err := db.CreateTable(&models.Device{}).Error; err != nil {
		return err
	}
	if err := db.CreateTable(&models.RefreshToken{}).Error; err != nil {
		return err
	}
	if err := db.CreateTable(&models.Salt{}).Error; err != nil {
		return err
	}
	if err := db.CreateTable(&models.Service{}).Error; err != nil {
		return err
	}
	return nil
}
