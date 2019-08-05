package config

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
