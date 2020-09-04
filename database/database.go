package database

import (
	// database driver

	"database/sql"
	"fmt"

	"github.com/MuShare/pluto/config"
)

func NewDatabase(config *config.Config) (*sql.DB, error) {
	dbcfg := config.Database
	connect, err := generateConnectionSchema(dbcfg)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("mysql", connect)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(dbcfg.MaxOpenConns)
	db.SetMaxIdleConns(dbcfg.MaxIdleConns)
	return db, nil
}

func generateConnectionSchema(dbcfg *config.DatabaseConfig) (string, error) {
	switch dbcfg.Type.String() {
	case "mysql":
		return dbcfg.User + ":" + dbcfg.Password + "@tcp(" + dbcfg.Host + ":" + dbcfg.Port.String() + ")/" + dbcfg.DB +
			"?charset=utf8mb4&parseTime=True&loc=Local", nil
	}
	return "", fmt.Errorf("%s db type is not support", dbcfg.Type.String())
}
