package mail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"

	"github.com/leeif/pluto/config"
)

type Mail struct {
	SMTP     string
	User     string
	Password string
}

func (m *Mail) Send(recv, subj, contentType, body string) error {

	from := mail.Address{"", m.User}
	to := mail.Address{"", recv}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "Content-Type: " + contentType + "; charset=UTF-8\r\n" + body

	// Connect to the SMTP Server
	host, _, err := net.SplitHostPort(m.SMTP)

	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", m.User, m.Password, host)

	// TLS config
	tlsconfig := &tls.Config{
		// InsecureSkipVerify: true,
		ServerName: host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", m.SMTP, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	err = c.Quit()
	if err != nil {
		return err
	}
	return nil
}

func NewMail() *Mail {
	c := config.GetConfig().Mail
	if c.SMTP.String() == "" {
		return nil
	}
	mail := &Mail{
		SMTP:     c.SMTP.String(),
		User:     *c.User,
		Password: *c.Password,
	}
	return mail
}
