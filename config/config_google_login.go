package config

type GoogleLoginConfig struct {
	Aud *string `kiper_value:"name:aud;help:audience"`
}

func newGoogleLoginConfig() *GoogleLoginConfig {
	return &GoogleLoginConfig{}
}
