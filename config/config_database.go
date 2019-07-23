package config

import "errors"

type DatabaseConfig struct {
	Type     *DBType     `pluto_value:"path"`
	Host     *BaseString `pluto_value:"host"`
	User     *BaseString `pluto_value:"name"`
	Password *BaseString `pluto_value:"password"`
	Port     *Port       `pluto_value:"port"`
	DB       *BaseString `pluto_value:"db"`
}

type DBType struct {
	t string
}

func (dt *DBType) Set(t string) error {
	switch t {
	case "mysql":
		dt.t = t
	default:
		return errors.New("Database type is not supported")
	}
	return nil
}

func (dt *DBType) String() string {
	return dt.t
}

func newDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Type:     &DBType{},
		Host:     &BaseString{},
		User:     &BaseString{},
		Password: &BaseString{},
		Port:     &Port{},
		DB:       &BaseString{},
	}
}
