package config

import "github.com/pkg/errors"

type MailConfig struct {
	SMTP     *SMTP   `kiper_value:"name:smtp"`
	User     *string `kiper_value:"name:user"`
	Password *string `kiper_value:"name:password"`
}

type SMTP struct {
	s string
}

func (smtp *SMTP) Set(s string) error {
	if s == "" {
		return errors.New("smtp server can not be empty")
	}
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
