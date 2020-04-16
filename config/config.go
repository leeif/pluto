package config

import (
	"fmt"
	"path/filepath"

	"log"

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
	Token       *TokenConfig       `kiper_config:"name:token"`
	View        *ViewConfig        `kiper_config:"name:view"`
	Admin       *AdminConfig       `kiper_config:"name:admin"`
	Cors        *CorsConfig        `kiper_config:"name:cors"`
	Registry    *RegistryConfig    `kiper_config:"name:registry"`
	OAuth       *OAuthConfig       `kiper_config:"name:oauth"`
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
		Token:       newTokenConfig(),
		View:        newViewConfig(),
		Admin:       newAdminConfig(),
		Cors:        newCorsConfig(),
		Registry:    newRegistryConfig(),
		OAuth:       newOAuthConfig(),
	}
	name := ""
	ag := make([]string, 0)
	if len(args) > 0 {
		name = filepath.Base(args[0])
		ag = args[1:]
	}

	kiper := kiper.NewKiper(name, "Pluto server")
	kiper.Kingpin.Version(version)
	kiper.Kingpin.HelpFlag.Short('h')

	kiper.SetConfigFileFlag("config.file", "config file", "./config.json")

	if err := kiper.Parse(c, ag); err != nil {
		return nil, err
	}
	c.Version = version

	printConfig(c)

	return c, nil
}

func printConfig(config *Config) {
	log.Println(fmt.Sprintf("AccessTokenExpire: %d", config.Token.AccessTokenExpire))
	log.Println(fmt.Sprintf("RegisterVerifyTokenExpire: %d", config.Token.RegisterVerifyTokenExpire))
	log.Println(fmt.Sprintf("ResetPasswordTokenExpire: %d", config.Token.ResetPasswordTokenExpire))
	log.Println(fmt.Sprintf("RefeshTokenExpire: %d", config.Token.RefreshTokenExpire))
}
