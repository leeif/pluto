package mail

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/MuShare/pluto/utils/view"
	"github.com/nicksnyder/go-i18n/v2/i18n"

	"github.com/MuShare/pluto/config"
	"github.com/MuShare/pluto/utils/jwt"

	perror "github.com/MuShare/pluto/datatype/pluto_error"
)

type Mail struct {
	config *config.Config
	bundle *i18n.Bundle
}

type SendMailRequest struct {
	To          string `json:"to"`
	Subject     string `json:"subject"`
	ContentType string `json:"content_type"`
	Body        string `json:"body"`
}

type SendMailResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (m *Mail) Send(recv, subj, contentType, body string) error {
	requestJson, err := json.Marshal(SendMailRequest{
		To:          recv,
		Subject:     subj,
		ContentType: contentType,
		Body:        body,
	})
	if err != nil {
		return err
	}
	res, err := http.Post(fmt.Sprintf("%s/api/v1/send-mail", *m.config.Mail.MailSenderPoolBaseUrl), "application/json", bytes.NewBuffer(requestJson))
	if err != nil {
		return err
	}

	defer res.Body.Close()
	var response SendMailResponse
	if res.StatusCode == http.StatusOK {
		if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
			return err
		}
		if response.Code != 200 {
			return errors.New(fmt.Sprintf("mail sender pool response is %d", response.Code))
		}
	} else {
		return errors.New(fmt.Sprintf("mail sender pool http request failed with status: %d", res.StatusCode))
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

// appName 是字符串名称，用于邮件模板
func (m *Mail) SendResetPassword(appID, address string, baseURL string, userLanguage string, appName string) *perror.PlutoError {
	prp := jwt.NewPasswordResetPayload(appID, address, m.config.Token.ResetPasswordTokenExpire)
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
	if *c.MailSenderPoolBaseUrl == "" {
		return nil, perror.ServerError.Wrapper(errors.New("mail sender pool base url is not set"))
	}
	mail := &Mail{
		config: config,
		bundle: bundle,
	}
	return mail, nil
}
