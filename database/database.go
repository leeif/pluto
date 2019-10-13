package database

import (
	// database driver

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	// gorm mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/leeif/pluto/config"
)

var database *gorm.DB

func NewDatabase(config *config.Config) (*gorm.DB, error) {
	if database == nil {
		dbcfg := config.Database
		con := generateConnectionSchema(dbcfg)
		db, err := gorm.Open(dbcfg.Type.String(), con)
		if err != nil {
			return nil, err
		}
		db.DB().SetMaxIdleConns(10)
		// db.LogMode(false)
		database = db
	}
	return database, nil
}

func generateConnectionSchema(dbcfg *config.DatabaseConfig) string {
	switch dbcfg.Type.String() {
	case "mysql":
		return dbcfg.User + ":" + dbcfg.Password + "@tcp(" + dbcfg.Host + ":" + dbcfg.Port.String() + ")/" + dbcfg.DB +
			"?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return ""
}
