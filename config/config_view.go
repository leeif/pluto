package config

type ViewConfig struct {
	Path      string `kiper_value:"name:path;help:path of html view files;default:./views"`
	Languages string `kiper_value:"name:language;help:support languages;default:en,zh,ja"`
}

func newViewConfig() *ViewConfig {
	return &ViewConfig{}
}
