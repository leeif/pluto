package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/leeif/pluto/utils/jwt"

	"github.com/leeif/pluto/datatype/request"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type errorResponse struct {
	Status string `json:"status"`
	Error  Error  `json:"error"`
}

type okResponse struct {
	Status string                 `json:"status"`
	Body   map[string]interface{} `json:"body"`
}

func sameResponse(origin interface{}, expect interface{}) (bool, error) {
	o, err := json.Marshal(origin)
	if err != nil {
		return false, err
	}
	e, err := json.Marshal(expect)
	if err != nil {
		return false, err
	}
	originStr := strings.ReplaceAll(string(o), " ", "")
	expectStr := strings.ReplaceAll(string(e), " ", "")
	if originStr != expectStr {
		return false, nil
	}
	return true, nil
}

func testMailRegisterBadRequest() error {
	url := "http://localhost:8010/api/user/register"
	payload := request.MailRegister{
		Mail:     "geeklyf@hotmail.com",
		Password: "test",
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp, err := http.Post(url, "application/json", bytes.NewReader(b)); err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	} else {
		if resp.StatusCode != http.StatusBadRequest {
			return fmt.Errorf("%v request expect bad request (400) status code, but %v", url, resp.StatusCode)
		}
	}
	return nil
}

func testMailRegisterOK() error {
	url := "http://localhost:8010/api/user/register"
	payload := request.MailRegister{
		Mail:     "test@gmail.com",
		Password: "test",
		Name:     "leeif",
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(b))

	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v request expect ok (200) status code, but %v", url, resp.StatusCode)
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}
	origin := okResponse{}
	err = json.Unmarshal(b, &origin)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	expect := okResponse{
		Status: "ok",
		Body: map[string]interface{}{
			"mail":     "test@gmail.com",
			"verified": true,
		},
	}

	same, err := sameResponse(origin, expect)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if !same {
		return fmt.Errorf("Expect response %v, but %v", expect, origin)
	}

	return nil
}

func testRegisterVerifyOK() error {
	rvp := jwt.NewRegisterVerifyPayload(1, 60*60)
	token, perror := jwt.GenerateRSAJWT(rvp)

	if perror != nil {
		return fmt.Errorf("Expect no error, but %v", perror.LogError)
	}

	url := "http://localhost:8010/mail/verify/" + token.B64String()
	resp, err := http.Get(url)

	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Expect ok (200) status code, but %v", resp.StatusCode)
	}

	return nil
}

var accessToken = ""
var refreshToken = ""

func testMailLoginOK() error {
	url := "http://localhost:8010/api/user/login"
	payload := request.MailLogin{
		Mail:     "test@gmail.com",
		Password: "test",
		AppID:    "test",
		DeviceID: "xxx",
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(b))

	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v request expect ok (200) status code, but %v", url, resp.StatusCode)
	}

	b, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}
	origin := okResponse{}
	err = json.Unmarshal(b, &origin)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}
	if jwt, ok := origin.Body["jwt"]; ok {
		accessToken = jwt.(string)
	} else {
		return fmt.Errorf("Expect to contain jwt filed in body")
	}
	if rt, ok := origin.Body["refresh_token"]; ok {
		refreshToken = rt.(string)
	} else {
		return fmt.Errorf("Expect to contain refresh_token filed in body")
	}
	return nil
}

func testGetUserInfo() error {
	url := "http://localhost:8010/api/user/info/me"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}
	req.Header.Set("Authorization", "jwt "+accessToken)
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v request expect ok (200) status code, but %v", url, resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}
	origin := okResponse{}
	err = json.Unmarshal(b, &origin)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	log.Printf("avatar url: %v", origin.Body["avatar"])

	resp, err = http.Get(origin.Body["avatar"].(string))
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v request expect ok (200) status code, but %v", origin.Body["avatar"], resp.StatusCode)
	}

	return nil
}

func testUpdateUserInfo() error {
	url := "http://localhost:8010/api/user/info/me/update"
	payload := request.UpdateUserInfo{
		Name: "test update",
	}

	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	fmt.Println(string(b))

	req, err := http.NewRequest("POST", url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	req.Header.Set("Authorization", "jwt "+accessToken)
	req.Header.Set("Content-type", "application/json")
	client := http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v request expect ok (200) status code, but %v", url, resp.StatusCode)
	}

	return nil
}
