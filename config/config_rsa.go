package config

type RSAConfig struct {
	Path *FilePath `pluto_value:"path"`
	Name *FileName `pluto_value:"name"`
}

type FilePath struct {
	s string
}

func (p *FilePath) Set(s string) error {
	p.s = s
	return nil
}

func (p *FilePath) String() string {
	return p.s
}

type FileName struct {
	s string
}

func (n *FileName) Set(s string) error {
	n.s = s
	return nil
}

func (n *FileName) String() string {
	return n.s
}

func newRSAConfig() *RSAConfig {
	return &RSAConfig{
		Name: &FileName{},
		Path: &FilePath{},
	}
}
