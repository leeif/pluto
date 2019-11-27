package config

type AdminConfig struct {
	Application string `kiper_value:"name:application;help:default admin application for pluto;default:pluto"`
	Role        string `kiper_value:"name:user;help:default role of admin application;default:root"`
	Scope       string `kiper_value:"name:scope;help:default role of admin application;default:admin"`
}

func newAdminConfig() *AdminConfig {
	return &AdminConfig{}
}
