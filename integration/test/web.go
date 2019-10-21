package test

import (
	"fmt"
	"net/http"

	"github.com/leeif/pluto/utils/jwt"
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
	prrp := jwt.NewPasswordResetResultPayload(true, 10*60)
	token, perror := jwt.GenerateRSAJWT(prrp)

	if perror != nil {
		return fmt.Errorf("Expect no error, but %v", perror.LogError)
	}

	url := "http://localhost:8010/password/reset/result/" + token.B64String()
	fmt.Println(url)
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
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}
	c := http.Client{}
	resp, err := c.Do(req)

	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusInternalServerError {
		return fmt.Errorf("Expect internal server error (500) status code, but %v", resp.StatusCode)
	}
	return nil
}

func testPasswordResetOK() error {
	prp := jwt.NewPasswordResetPayload("test@gmail.com", 60*60)
	token, perror := jwt.GenerateRSAJWT(prp)

	if perror != nil {
		return fmt.Errorf("Expect no error, but %v", perror.LogError)
	}

	url := "http://localhost:8010/password/reset/" + token.B64String()
	fmt.Println(url)
	resp, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expect ok (200) status code, but %v", resp.StatusCode)
	}
	return nil
}
