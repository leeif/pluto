package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/leeif/pluto/datatype/request"
)

func testGetPublicKey() error {
	url := "http://localhost:8010/api/auth/publickey"
	resp, err := http.Get(url)
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
	publicKey := origin.Body["public_key"]
	f, err := os.Open(path.Join(rsaDir, "id_rsa_test.pub"))
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}
	b, err = ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}

	if publicKey != string(b) {
		return fmt.Errorf("public key not equal to the loacl file")
	}
	return nil
}

func testRefreshAccessToken() error {
	url := "http://localhost:8010/api/auth/refresh"
	payload := request.RefreshAccessToken{
		RefreshToken: refreshToken,
		UseID:        1,
		AppID:        "xxx",
		DeviceID:     "xxx",
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
	return nil
}
