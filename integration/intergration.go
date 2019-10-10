package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/utils/rsa"
)

type testCase struct {
	Name string
	Func func() error
}

var testCases = []testCase{
	{
		Name: "testHealthCheck",
		Func: testHealthCheck,
	},
	{
		Name: "initRSA",
		Func: initRSA,
	},
	{
		Name: "testMailRegisterBadRequest",
		Func: testMailRegisterBadRequest,
	},
	{
		Name: "testMailRegisterOK",
		Func: testMailRegisterOK,
	},
	{
		Name: "testRegisterVerifyOK",
		Func: testRegisterVerifyOK,
	},
	{
		Name: "testMailLoginOK",
		Func: testMailLoginOK,
	},
	{
		Name: "testGetUserInfo",
		Func: testGetUserInfo,
	},
	{
		Name: "testGetPublicKey",
		Func: testGetPublicKey,
	},
	{
		Name: "testRefreshAccessToken",
		Func: testRefreshAccessToken,
	},
	{
		Name: "testResetPassword",
		Func: testResetPassword,
	},
	{
		Name: "testNewPasswordLogin",
		Func: testNewPasswordLogin,
	},
	{
		Name: "testPasswordResetResultFail",
		Func: testPasswordResetResultFail,
	},
	{
		Name: "testPasswordResetResultOK",
		Func: testPasswordResetResultOK,
	},
	{
		Name: "testPasswordResetFail",
		Func: testPasswordResetFail,
	},
	{
		Name: "testPasswordResetOK",
		Func: testPasswordResetOK,
	},
}

func testHealthCheck() error {
	url := "http://localhost:8010/healthcheck"
	for i := 0; i < 100; i++ {
		log.Printf("retry count: %v\n", i)
		resp, err := http.Get(url)
		time.Sleep(time.Duration(5) * time.Second)
		if err != nil {
			continue
		}
		if resp.StatusCode == http.StatusOK {
			return nil
		}
	}
	return errors.New("Healthcheck failed")
}

func initRSA() error {
	cfg := config.Config{}
	cfg.RSA = &config.RSAConfig{}
	name := "ids_rsa_test"
	cfg.RSA.Name = &name
	path := "./docker"
	cfg.RSA.Path = &path
	if err := rsa.Init(&cfg); err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}
	return nil
}

func main() {
	// test cases
	for _, tc := range testCases {
		log.Printf("====== start %v ======\n", tc.Name)
		err := tc.Func()
		time.Sleep(time.Duration(1) * time.Second)
		if err != nil {
			log.Panicf("Error: %v", err)
		}
		log.Printf("test ok\n")
		log.Printf("====== end %v ======\n\n", tc.Name)
	}
}
