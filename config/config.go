package config

import "github.com/pkg/errors"

var config *Config

type PlutoConfig interface {
}

type PlutoValue interface {
	Set(string) error
	String() string
}

type Config struct {
	ConfigFile *string         `kiper_value:"name:config.file;default:./config.json"`
	Server     *ServerConfig   `kiper_config:"name:server"`
	Log        *LogConfig      `kiper_config:"name:log"`
	RSA        *RSAConfig      `kiper_config:"name:rsa"`
	Database   *DatabaseConfig `kiper_config:"name:database"`
	Mail       *MailConfig     `kiper_config:"name:mail"`
}

func CheckConfig(c *Config) error {
	if c.Mail.SMTP.String() == "" {
		return errors.New("smtp can not be empty")
	}
	return nil
}

func GetConfig() *Config {
	if config == nil {
		config = NewConfig()
	}
	return config
}

func NewConfig() *Config {
	c := &Config{
		Log:      newLogConfig(),
		Server:   newServerConfig(),
		RSA:      newRSAConfig(),
		Database: newDatabaseConfig(),
		Mail:     newMailConfig(),
	}
	return c
}
