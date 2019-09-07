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

	b64 "encoding/base64"

	"github.com/alecthomas/template"
	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/utils/jwt"

	perror "github.com/leeif/pluto/datatype/pluto_error"
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

func SendRegisterVerify(userID uint, mail string) *perror.PlutoError {
	// expire time 10 mins
	token, err := jwt.GenerateJWT(jwt.Head{Type: jwt.REGISTERVERIFY}, &jwt.RegisterVerifyPayload{UserID: userID}, 10*60)
	if err != nil {
		return err.Wrapper(errors.New("JWT token generate failed"))
	}

	if m := NewMail(); m != nil {
		dir, _ := os.Getwd()
		t := template.Must(template.ParseFiles(path.Join(dir, "views", "register_verify_mail.html")))
		var buffer bytes.Buffer
		type Data struct {
			BaseURL string
			Token   string
		}
		baseURL := config.GetConfig().Server.BaseURL
		t.Execute(&buffer, Data{Token: b64.StdEncoding.EncodeToString([]byte(token)), BaseURL: *baseURL})
		if err := m.Send(mail, "[MuShare]Mail Verification", "text/html", buffer.String()); err != nil {
			return perror.ServerError.Wrapper(errors.New("Mail sending failed: " + err.Error()))
		}
	} else {
		return perror.ServerError.Wrapper(errors.New("Mail sender is not defined"))
	}
	return nil
}

func SendResetPassword(mail string) *perror.PlutoError {
	// expire time 10 mins
	token, err := jwt.GenerateJWT(jwt.Head{Type: jwt.PASSWORDRESET}, &jwt.PasswordResetPayload{Mail: mail}, 10*60)
	if err != nil {
		return err.Wrapper(errors.New("JWT token generate failed"))
	}

	if m := NewMail(); m != nil {
		dir, _ := os.Getwd()
		t := template.Must(template.ParseFiles(path.Join(dir, "views", "password_reset_mail.html")))
		var buffer bytes.Buffer
		type Data struct {
			BaseURL string
			Token   string
		}
		baseURL := config.GetConfig().Server.BaseURL
		t.Execute(&buffer, Data{Token: b64.StdEncoding.EncodeToString([]byte(token)), BaseURL: *baseURL})
		if err := m.Send(mail, "[MuShare]Password Reset", "text/html", buffer.String()); err != nil {
			return perror.ServerError.Wrapper(errors.New("Mail sending failed: " + err.Error()))
		}
	} else {
		return perror.ServerError.Wrapper(errors.New("Mail sender is not defined"))
	}
	return nil
}
