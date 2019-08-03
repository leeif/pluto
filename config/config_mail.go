package config

type MailConfig struct {
	SMTP     *string `pluto_value:"name:smtp"`
	User     *string `pluto_value:"name:user"`
	Password *string `pluto_value:"name:format"`
}

func newMailConfig() *MailConfig {
	return &MailConfig{}
}
