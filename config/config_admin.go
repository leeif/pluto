package config

type AdminConfig struct {
	Mail string `kiper_value:"name:mail;help:default mail register for pluto admin"`
	Name string `kiper_value:"name:name;help:default user name of pluto admin;default:root"`
}

func newAdminConfig() *AdminConfig {
	return &AdminConfig{}
}
