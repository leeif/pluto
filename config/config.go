package config

import (
	"path/filepath"

	"github.com/leeif/kiper"
)

var config *Config

type PlutoConfig interface {
}

type PlutoValue interface {
	Set(string) error
	String() string
}

type Config struct {
	Version     string
	Server      *ServerConfig      `kiper_config:"name:server"`
	Log         *LogConfig         `kiper_config:"name:log"`
	RSA         *RSAConfig         `kiper_config:"name:rsa"`
	Database    *DatabaseConfig    `kiper_config:"name:database"`
	Mail        *MailConfig        `kiper_config:"name:mail"`
	Avatar      *AvatarConfig      `kiper_config:"name:avatar"`
	GoogleLogin *GoogleLoginConfig `kiper_config:"name:google_login"`
	WechatLogin *WechatLoginConfig `kiper_config:"name:wechat_login"`
	AppleLogin  *AppleLoginConfig  `kiper_config:"name:apple_login"`
	JWT         *JWTConfig         `kiper_config:"name:jwt"`
	View        *ViewConfig        `kiper_config:"name:view"`
	Admin       *AdminConfig       `kiper_config:"name:admin"`
}

func NewConfig(args []string, version string) (*Config, error) {
	c := &Config{
		Log:         newLogConfig(),
		Server:      newServerConfig(),
		RSA:         newRSAConfig(),
		Database:    newDatabaseConfig(),
		Mail:        newMailConfig(),
		Avatar:      newAvatarConfig(),
		GoogleLogin: newGoogleLoginConfig(),
		WechatLogin: newWechatLoginConfig(),
		AppleLogin:  newAppleLoginConfig(),
		JWT:         newJWTConfig(),
		View:        newViewConfig(),
		Admin:       newAdminConfig(),
	}
	kiper := kiper.NewKiper(filepath.Base(args[0]), "Pluto server")
	kiper.Kingpin.Version(version)
	kiper.Kingpin.HelpFlag.Short('h')

	kiper.SetConfigFileFlag("config.file", "config file", "./config.json")

	if err := kiper.Parse(c, args[1:]); err != nil {
		return nil, err
	}
	c.Version = version
	return c, nil
}
