package config

type CrosConfig struct {
	AllowedOrigins []string `kiper_value:"name:allow_origins;help:allow origins url;default:*"`
}

func newCrosConfig() *CrosConfig {
	return &CrosConfig{}
}
