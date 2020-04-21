package config

type AdminConfig struct {
	Mail string `kiper_value:"name:mail;help:default mail register for pluto admin"`
}

func newAdminConfig() *AdminConfig {
	return &AdminConfig{}
}
