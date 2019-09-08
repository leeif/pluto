package main

import (
	"errors"
	"log"
	"net/http"
	"os/exec"
	"time"
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
	for i := 0; i < 500; i++ {
		log.Printf("try count: %v\n", i)
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

func main() {
	log.Println("docker-compose -f docker/docker-compose.yml up -d")
	cmd := exec.Command("docker-compose", "-f", "docker/docker-compose.yml", "up", "-d")
	cmd.Start()
	time.Sleep(time.Duration(10) * time.Second)
	defer func() {
		log.Println("docker-compose -f docker/docker-compose.yml down --rmi all")
		cmd := exec.Command("docker-compose", "-f", "docker/docker-compose.yml", "down", "--rmi", "all")
		cmd.Start()
		time.Sleep(time.Duration(10) * time.Second)
	}()
	for _, tc := range testCases {
		log.Printf("====== start %v ======\n", tc.Name)
		err := tc.Func()
		if err != nil {
			log.Panicf("Error: %v", err)
		}
		log.Printf("test ok\n")
		log.Printf("====== end %v ======\n\n", tc.Name)
	}
}
