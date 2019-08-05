package config

type MailConfig struct {
	SMTP     *string `kiper_value:"name:smtp"`
	User     *string `kiper_value:"name:user"`
	Password *string `kiper_value:"name:password"`
}

func newMailConfig() *MailConfig {
	return &MailConfig{}
}
