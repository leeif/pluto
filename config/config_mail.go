package config

type MailConfig struct {
	SMTP     *SMTP   `kiper_value:"name:smtp"`
	User     *string `kiper_value:"name:user"`
	Password *string `kiper_value:"name:password"`
}

type SMTP struct {
	s string
}

func (smtp *SMTP) Set(s string) error {
	smtp.s = s
	return nil
}

func (smtp *SMTP) String() string {
	return smtp.s
}

func newMailConfig() *MailConfig {
	return &MailConfig{
		SMTP: &SMTP{},
	}
}
