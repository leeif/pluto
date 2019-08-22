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

func GetDatabase() (*gorm.DB, error) {
	if database == nil {
		config := config.GetConfig().Database
		con := generateConnectionSchema(config)
		db, err := gorm.Open(config.Type.String(), con)
		if err != nil {
			return nil, err
		}
		db.DB().SetMaxIdleConns(10)
		// db.LogMode(false)
		database = db
	}
	return database, nil
}

func generateConnectionSchema(config *config.DatabaseConfig) string {
	switch config.Type.String() {
	case "mysql":
		return *config.User + ":" + *config.Password + "@tcp(" + *config.Host + ":" + config.Port.String() + ")/" + *config.DB +
			"?charset=utf8&parseTime=True&loc=Local"
	}
	return ""
}
