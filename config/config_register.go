package config

type RegisterConfig struct {
	Consul        bool   `kiper_value:"name:consul;help:use consul register;default:false"`
	ConsulAddress string `kiper_value:"name:consul_address;help:consul address;default:127.0.0.1"`
	ConsulPort    *Port  `kiper_value:"name:consul_port;help:consul port;default:8500"`
	ServiceName   string `kiper_value:"name:service_name;help:register service name;default:pluto"`
}

func newRegisterConfig() *RegisterConfig {
	return &RegisterConfig{
		ConsulPort: &Port{},
	}
}
