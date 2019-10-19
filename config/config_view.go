package config

type ViewConfig struct {
	Path string `kiper_value:"name:path;help:path of html view files;default:./views"`
}

func newViewConfig() *ViewConfig {
	return &ViewConfig{}
}
