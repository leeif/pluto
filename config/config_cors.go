package config

type CorsConfig struct {
	AllowedOrigins []string `kiper_value:"name:allow_origins;help:allow origins url;default:*"`
	AllowedHeaders []string `kiper_value:"name:allow_headers;help:allow headers;default:*"`
}

func newCorsConfig() *CorsConfig {
	return &CorsConfig{}
}
