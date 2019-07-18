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
	Server     *ServerConfig `pluto_config:"server"`
	Log        *LogConfig    `pluto_config:"log"`
	RSA        *RSAConfig    `pluto_config:"rsa"`
}

func (c *Config) setFlag(a *kingpin.Application, args []string) error {
	a.Flag("config.file", "configure file path").Default("./config.json").StringVar(&c.ConfigFile)
	c.setPlutoServerFlag(a)
	c.setLogFlag(a)
	c.setRSAFlag(a)
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
}

func (c *Config) setRSAFlag(a *kingpin.Application) {
	if c.RSA.Name != nil {
		a.Flag("rsa.name", "rsa public/private key path").Default("ids_ra").SetValue(c.RSA.Name)
	}

	if c.RSA.Path != nil {
		a.Flag("rsa.path", "log format: json, logfmt").Default("./").SetValue(c.RSA.Path)
	}
}

func (c *Config) loadConfigFile() error {
	viper.SetConfigFile("/etc/pluto/config.json")
	viper.SetConfigFile("$HOME/.pluto/config.json")
	viper.SetConfigFile("./config.json")
	viper.SetConfigFile(c.ConfigFile)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
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
	mergeCommandLineWithConfigFile(c, viper.AllSettings())
}

func GetConfig() *Config {
	if config == nil {
		config = &Config{
			Log:    newLogConfig(),
			Server: newServerConfig(),
			RSA:    newRSAConfig(),
		}
	}
	return config
}
