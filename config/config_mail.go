package config

type MailConfig struct {
	MailSenderPoolBaseUrl *string `kiper_value:"name:mail-sender-pool-base-url;help:mail sender pool host"`
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
	return &MailConfig{}
}
