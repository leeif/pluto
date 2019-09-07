package config

import (
	"os"
	"path/filepath"

	"github.com/leeif/kiper"
	"github.com/pkg/errors"
)

var config *Config

type PlutoConfig interface {
}

type PlutoValue interface {
	Set(string) error
	String() string
}

type Config struct {
	ConfigFile  *string            `kiper_value:"name:config.file;default:./config.json"`
	Server      *ServerConfig      `kiper_config:"name:server"`
	Log         *LogConfig         `kiper_config:"name:log"`
	RSA         *RSAConfig         `kiper_config:"name:rsa"`
	Database    *DatabaseConfig    `kiper_config:"name:database"`
	Mail        *MailConfig        `kiper_config:"name:mail"`
	Avatar      *AvatarConfig      `kiper_config:"name:avatar"`
	GoogleLogin *GoogleLoginConfig `kiper_config:"name:google_login"`
}

func (c *Config) checkConfig() error {
	if c.Mail.SMTP.String() == "" {
		return errors.New("smtp can not be empty")
	}
	return nil
}

func NewConfig() (*Config, error) {
	c := &Config{
		Log:         newLogConfig(),
		Server:      newServerConfig(),
		RSA:         newRSAConfig(),
		Database:    newDatabaseConfig(),
		Mail:        newMailConfig(),
		Avatar:      newAvatarConfig(),
		GoogleLogin: newGoogleLoginConfig(),
	}
	kiper := kiper.NewKiper(filepath.Base(os.Args[0]), "Pluto server")
	kiper.GetKingpinInstance().HelpFlag.Short('h')

	if err := kiper.ParseCommandLine(c, os.Args[1:]); err != nil {
		return nil, err
	}

	if err := kiper.ParseConfigFile(*c.ConfigFile); err != nil {
		return nil, err
	}

	if err := kiper.MergeConfigFile(c); err != nil {
		return nil, err
	}

	if err := c.checkConfig(); err != nil {
		return nil, err
	}
	return c, nil
}
