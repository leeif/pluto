package mail

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"path"

	"github.com/alecthomas/template"
	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/utils/jwt"

	perror "github.com/leeif/pluto/datatype/pluto_error"
)

type Mail struct {
	config *config.Config
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

func (m *Mail) SendRegisterVerify(userID uint, address string, baseURL string) *perror.PlutoError {
	// expire time 10 mins
	token, err := jwt.GenerateJWT(jwt.Head{Type: jwt.REGISTERVERIFY}, &jwt.RegisterVerifyPayload{UserID: userID}, m.config.JWT.RegisterVerifyTokenExpire)
	if err != nil {
		return err.Wrapper(errors.New("JWT token generate failed"))
	}

	dir, _ := os.Getwd()
	t := template.Must(template.ParseFiles(path.Join(dir, "views", "register_verify_mail.html")))
	var buffer bytes.Buffer
	type Data struct {
		BaseURL string
		Token   string
	}
	t.Execute(&buffer, Data{Token: token.B64String(), BaseURL: baseURL})
	if err := m.Send(address, "[MuShare]Mail Verification", "text/html", buffer.String()); err != nil {
		return perror.ServerError.Wrapper(errors.New("Mail sending failed: " + err.Error()))
	}

	return nil
}

func (m *Mail) SendResetPassword(address string, baseURL string) *perror.PlutoError {
	// expire time 10 mins
	token, err := jwt.GenerateJWT(jwt.Head{Type: jwt.PASSWORDRESET}, &jwt.PasswordResetPayload{Mail: address}, m.config.JWT.ResetPasswordTokenExpire)
	if err != nil {
		return err.Wrapper(errors.New("JWT token generate failed"))
	}

	dir, _ := os.Getwd()
	t := template.Must(template.ParseFiles(path.Join(dir, "views", "password_reset_mail.html")))
	var buffer bytes.Buffer
	type Data struct {
		BaseURL string
		Token   string
	}
	t.Execute(&buffer, Data{Token: token.B64String(), BaseURL: baseURL})
	if err := m.Send(address, "[MuShare]Password Reset", "text/html", buffer.String()); err != nil {
		return perror.ServerError.Wrapper(errors.New("Mail sending failed: " + err.Error()))
	}
	return nil
}

func NewMail(config *config.Config) (*Mail, *perror.PlutoError) {
	c := config.Mail
	if c.SMTP.String() == "" {
		return nil, perror.ServerError.Wrapper(errors.New("smtp server is not set"))
	}
	mail := &Mail{
		config: config,
	}
	return mail, nil
}
