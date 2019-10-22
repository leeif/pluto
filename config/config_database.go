package config

import "errors"

type DatabaseConfig struct {
	Type     *DBType `kiper_value:"name:type;help:database type;default:mysql"`
	Host     string  `kiper_value:"name:host;help:database host;default:127.0.0.1"`
	User     string  `kiper_value:"name:user;help:database user;default:root"`
	Password string  `kiper_value:"name:password;help:database password"`
	Port     *Port   `kiper_value:"name:port;help:database port;default:3306"`
	DB       string  `kiper_value:"name:db;help:db name;default:pluto"`
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
		Type: &DBType{},
		Port: &Port{},
	}
}
