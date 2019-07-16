package config

type ServerConfig struct {
	Port *Port `pluto_value:"port"`
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
