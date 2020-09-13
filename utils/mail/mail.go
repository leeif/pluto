package mail

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net"
	"net/mail"
	"net/smtp"

	"github.com/MuShare/pluto/utils/view"

	"github.com/MuShare/pluto/config"
	"github.com/MuShare/pluto/utils/jwt"

	perror "github.com/MuShare/pluto/datatype/pluto_error"
)

type Mail struct {
	config *config.Config
	bundle *i18n.Bundle
}

func (m *Mail) Send(recv, subj, contentType, body string) error {

	from := mail.Address{"", *m.config.Mail.User}
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
	host, _, err := net.SplitHostPort(m.config.Mail.SMTP.String())

	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", *m.config.Mail.User, *m.config.Mail.Password, host)

	// TLS config
	tlsconfig := &tls.Config{
		// InsecureSkipVerify: true,
		ServerName: host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", m.config.Mail.SMTP.String(), tlsconfig)
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

func (m *Mail) SendPlainText(address, subject, text string) *perror.PlutoError {
	if err := m.Send(address, subject, "text/plain", text); err != nil {
		return perror.ServerError.Wrapper(errors.New("Mail sending failed: " + err.Error()))
	}

	return nil
}

func (m *Mail) SendRegisterVerify(userID uint, address string, baseURL string, language string, appName string) *perror.PlutoError {
	rvp := jwt.NewRegisterVerifyPayload(userID, m.config.Token.RegisterVerifyTokenExpire)
	token, perr := jwt.GenerateRSA256JWT(rvp)
	if perr != nil {
		return perr.Wrapper(errors.New("JWT token generate failed"))
	}

	vw, err := view.GetView()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	t, err := vw.Parse(language, "register_verify_mail.html")
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}
	var buffer bytes.Buffer
	type Data struct {
		BaseURL string
		Token   string
	}
	t.Execute(&buffer, Data{Token: token.B64String(), BaseURL: baseURL})
	localizer := i18n.NewLocalizer(m.bundle, language)
	subject, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "VerifyMailSubject",
		TemplateData: map[string]interface{}{
			"AppName": appName,
		},
	})
	if err != nil {
		subject = "[Pluto] Mail Confirmation"
	}
	if err := m.Send(address, subject, "text/html", buffer.String()); err != nil {
		return perror.ServerError.Wrapper(errors.New("Mail sending failed: " + err.Error()))
	}

	return nil
}

func (m *Mail) SendResetPassword(address string, baseURL string, userLanguage string, appName string) *perror.PlutoError {
	prp := jwt.NewPasswordResetPayload(address, m.config.Token.ResetPasswordTokenExpire)
	token, perr := jwt.GenerateRSA256JWT(prp)
	if perr != nil {
		return perr.Wrapper(errors.New("JWT token generate failed"))
	}

	vw, err := view.GetView()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	t, err := vw.Parse(userLanguage, "password_reset_mail.html")
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	var buffer bytes.Buffer
	type Data struct {
		BaseURL string
		Token   string
	}
	t.Execute(&buffer, Data{Token: token.B64String(), BaseURL: baseURL})
	localizer := i18n.NewLocalizer(m.bundle, userLanguage)
	subject, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: "ResetPasswordMailSubject",
		TemplateData: map[string]interface{}{
			"AppName": appName,
		},
	})
	if err != nil {
		subject = "[Pluto]Password Reset"
	}
	if err := m.Send(address, subject, "text/html", buffer.String()); err != nil {
		return perror.ServerError.Wrapper(errors.New("Mail sending failed: " + err.Error()))
	}
	return nil
}

func NewMail(config *config.Config, bundle *i18n.Bundle) (*Mail, *perror.PlutoError) {
	c := config.Mail
	if c.SMTP.String() == "" {
		return nil, perror.ServerError.Wrapper(errors.New("smtp server is not set"))
	}
	mail := &Mail{
		config: config,
		bundle: bundle,
	}
	return mail, nil
}
