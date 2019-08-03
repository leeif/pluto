package config

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var config *Config

type PlutoConfig interface {
}

type PlutoValue interface {
	Set(string) error
	String() string
}

type Config struct {
	ConfigFile string
	Server     *ServerConfig   `pluto_config:"server"`
	Log        *LogConfig      `pluto_config:"log"`
	RSA        *RSAConfig      `pluto_config:"rsa"`
	Database   *DatabaseConfig `pluto_config:"database"`
}

func (c *Config) setFlag(a *kingpin.Application, args []string) error {
	a.Flag("config.file", "configure file path").Default("./config.json").StringVar(&c.ConfigFile)
	c.setPlutoServerFlag(a)
	c.setLogFlag(a)
	c.setRSAFlag(a)
	c.setDatabaseFlag(a)
	_, err := a.Parse(args)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) setPlutoServerFlag(a *kingpin.Application) {
	if c.Server.Port != nil {
		a.Flag("server.port", "pluto server port").Default("8888").SetValue(c.Server.Port)
	}
}

func (c *Config) setLogFlag(a *kingpin.Application) {
	if c.Log.Level != nil {
		a.Flag("log.level", "log level: debug, info, warn, error").Default("info").SetValue(c.Log.Level)
	}

	if c.Log.Format != nil {
		a.Flag("log.format", "log format: json, logfmt").Default("logfmt").SetValue(c.Log.Format)
	}

	if c.Log.File != nil {
		a.Flag("log.file", "log file path").Default("").SetValue(c.Log.File)
	}
}

func (c *Config) setRSAFlag(a *kingpin.Application) {
	if c.RSA.Name != nil {
		a.Flag("rsa.name", "rsa public/private key name").Default("ids_ra").SetValue(c.RSA.Name)
	}

	if c.RSA.Path != nil {
		a.Flag("rsa.path", "rsa public/private key path").Default("./").SetValue(c.RSA.Path)
	}
}

func (c *Config) setDatabaseFlag(a *kingpin.Application) {
	if c.Database.Type != nil {
		a.Flag("database.type", "type of database").Default("mysql").SetValue(c.Database.Type)
	}

	if c.Database.Host != nil {
		a.Flag("database.host", "host of database").Default("127.0.0.1").SetValue(c.Database.Host)
	}

	if c.Database.User != nil {
		a.Flag("database.user", "user of database").Default("root").SetValue(c.Database.User)
	}

	if c.Database.Port != nil {
		a.Flag("database.port", "port of database").Default("3306").SetValue(c.Database.Port)
	}

	if c.Database.Password != nil {
		a.Flag("database.password", "password of database").Default("").SetValue(c.Database.Password)
	}

	if c.Database.DB != nil {
		a.Flag("database.db", "db of database").Default("").SetValue(c.Database.DB)
	}
}

func (c *Config) loadConfigFile() error {
	viper.SetConfigFile("/etc/pluto/config.json")
	viper.SetConfigFile("$HOME/.pluto/config.json")
	viper.SetConfigFile("./config.json")
	viper.SetConfigFile(c.ConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			// Config file not found; ignore error
		} else {
			// Config file was found but another error was produced
			return err
		}
	}

	return nil
}

func (c *Config) Parse(a *kingpin.Application, args []string) {
	if err := c.setFlag(a, args); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		a.Usage(args)
		os.Exit(2)
	}

	if err := c.loadConfigFile(); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing config file"))
		a.Usage(args)
		os.Exit(2)
	}

	// merge command line and config file settings
	if err := mergeCommandLineWithConfigFile(c, viper.AllSettings()); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error mergin config file"))
		os.Exit(2)
	}
}

func GetConfig() *Config {
	if config == nil {
		config = &Config{
			Log:      newLogConfig(),
			Server:   newServerConfig(),
			RSA:      newRSAConfig(),
			Database: newDatabaseConfig(),
		}
	}
	return config
}
