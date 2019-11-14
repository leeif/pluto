package config

type AppleLoginConfig struct {
	TeamID      string `kiper_value:"name:team_id;help:apple team id"`
	BundleID    string `kiper_value:"name:bundle_id;help:apple bundle id"`
	ClientID    string `kiper_value:"name:client_id;help:apple service id"`
	KeyID       string `kiper_value:"name:key_id;help:apple key id"`
	P8CertFile  string `kiper_value:"name:p8_cert_file;help:p8 cert file path"`
	RedirectURL string `kiper_value:"name:redirect_url;help:redirect url"`
}

func newAppleLoginConfig() *AppleLoginConfig {
	return &AppleLoginConfig{}
}
