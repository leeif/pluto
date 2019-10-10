package config

type ServerConfig struct {
	// skip sending verify mail when mail register
	SkipRegisterVerifyMail bool  `kiper_value:"name:skip_register_verify_mail;default:false"`
	Port                   *Port `kiper_value:"name:port;help:server listen port;default:8010"`
}

type Port struct {
	s string
}

func (p *Port) Set(s string) error {
	p.s = s
	return nil
}

func (p *Port) String() string {
	return p.s
}

func newServerConfig() *ServerConfig {
	return &ServerConfig{
		Port: &Port{},
	}
}
