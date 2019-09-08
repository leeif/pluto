package main

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/leeif/pluto/utils/rsa"
)

func testPasswordResetResultFail() error {

	url := "http://localhost:8010/password/reset/result/" + "random"
	resp, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		return fmt.Errorf("Expect internal server error (500) status code, but %v", resp.StatusCode)
	}
	return nil
}

func testPasswordResetResultOK() error {
	cfg := config.Config{}
	cfg.RSA = &config.RSAConfig{}
	name := "ids_rsa_test"
	cfg.RSA.Name = &name
	path := "./docker"
	cfg.RSA.Path = &path
	if err := rsa.Init(&cfg); err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}
	token, perror := jwt.GenerateJWT(jwt.Head{Type: jwt.PASSWORDRESETRESULT, Alg: jwt.ALGRAS},
		&jwt.PasswordResetResultPayload{Message: "success"}, 60*60)

	if perror != nil {
		return fmt.Errorf("Expect no error, but %v", perror.LogError)
	}

	url := "http://localhost:8010/password/reset/result/" + base64.StdEncoding.EncodeToString([]byte(token))
	resp, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expect ok (200) status code, but %v", resp.StatusCode)
	}
	return nil
}

func testPasswordResetFail() error {

	url := "http://localhost:8010/password/reset/" + "random"
	resp, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		return fmt.Errorf("Expect internal server error (500) status code, but %v", resp.StatusCode)
	}
	return nil
}

func testPasswordResetOK() error {
	cfg := config.Config{}
	cfg.RSA = &config.RSAConfig{}
	name := "ids_rsa_test"
	cfg.RSA.Name = &name
	path := "./docker"
	cfg.RSA.Path = &path
	if err := rsa.Init(&cfg); err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}
	token, perror := jwt.GenerateJWT(jwt.Head{Type: jwt.PASSWORDRESET, Alg: jwt.ALGRAS},
		&jwt.PasswordResetPayload{Mail: "test@gmail.com"}, 60*60)

	if perror != nil {
		return fmt.Errorf("Expect no error, but %v", perror.LogError)
	}

	url := "http://localhost:8010/password/reset/" + base64.StdEncoding.EncodeToString([]byte(token))
	resp, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expect ok (200) status code, but %v", resp.StatusCode)
	}
	return nil
}
