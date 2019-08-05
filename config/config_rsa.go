package config

type RSAConfig struct {
	Path *string `kiper_value:"name:path;help:rsa file path;default:./"`
	Name *string `kiper_value:"name:name;help:rsa file name;default:id_rsa"`
}

func newRSAConfig() *RSAConfig {
	return &RSAConfig{}
}
