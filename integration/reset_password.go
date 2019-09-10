package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/leeif/pluto/utils/jwt"

	"github.com/leeif/pluto/datatype/request"
)

func testResetPassword() error {
	token, perror := jwt.GenerateJWT(jwt.Head{Type: jwt.PASSWORDRESET, Alg: jwt.ALGRAS},
		&jwt.PasswordResetPayload{Mail: "test@gmail.com"}, 60*60)

	if perror != nil {
		return fmt.Errorf("Expect no error, but %v", perror.LogError)
	}

	url := "http://localhost:8010/api/user/password/reset"
	payload := request.ResetPassword{
		Token:    base64.StdEncoding.EncodeToString([]byte(token)),
		Password: "test_new",
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp, err := http.Post(url, "application/json", bytes.NewReader(b)); err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	} else {
		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("%v request expect ok (200) status code, but %v", url, resp.StatusCode)
		}
	}
	return nil
}

func testNewPasswordLogin() error {
	url := "http://localhost:8010/api/user/login"
	payload := request.MailLogin{
		Mail:     "test@gmail.com",
		Password: "test",
		DeviceID: "xxx",
		AppID:    "xxx",
	}

	// login with old password
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(b))

	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusForbidden {
		return fmt.Errorf("%v request expect forbidden (403) status code, but %v", url, resp.StatusCode)
	}

	// login with new password
	payload.Password = "test_new"
	b, err = json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	resp, err = http.Post(url, "application/json", bytes.NewReader(b))

	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v request expect ok (200) status code, but %v", url, resp.StatusCode)
	}

	return nil
}
