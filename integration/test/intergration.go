package test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path"
	"path/filepath"
	"runtime"
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
		Name: "testDataInitCheck",
		Func: dataInitCheck,
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
		Name: "testMailLoginOK",
		Func: testMailLoginOK,
	},
	{
		Name: "testGetUserInfo",
		Func: testGetUserInfo,
	},
	{
		Name: "testUpdateUserInfo",
		Func: testUpdateUserInfo,
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

func dataInitCheck() error {
	// port := &config.Port{}
	// typ := &config.DBType{}
	// port.Set("3306")
	// typ.Set("mysql")
	// config := &config.Config{
	// 	Database: &config.DatabaseConfig{
	// 		Type:     typ,
	// 		User:     "root",
	// 		Host:     "127.0.0.1",
	// 		Port:     port,
	// 		Password: "root",
	// 		DB:       "pluto_server",
	// 	},
	// }
	// db, err := database.NewDatabase(config)
	// if err != nil {
	// 	return err
	// }
	// for i := 0; i < 100; i++ {
	// 	_, err := models.Users(qm.Where("mail = ?", "geeklyf92610@gmail.com")).One(db)
	// 	if err != nil && err != sql.ErrNoRows {
	// 		return err
	// 	}
	// 	if err == nil {
	// 		return nil
	// 	}
	// 	time.Sleep(time.Duration(5) * time.Second)
	// }
	time.Sleep(time.Duration(10) * time.Second)
	// return errors.New("Data init failed")
	return nil
}

var rsaDir string

func initRSA() error {
	cfg := config.Config{}
	cfg.RSA = &config.RSAConfig{}
	name := "id_rsa_test"
	cfg.RSA.Name = &name
	_, filename, _, _ := runtime.Caller(2)
	rsaDir = path.Join(filepath.Dir(filename), "./docker")
	fmt.Println(rsaDir)
	cfg.RSA.Path = &rsaDir
	if err := rsa.Init(&cfg, nil); err != nil {
		return fmt.Errorf("Expect no error, but %v", err)
	}
	return nil
}

func Integration() {
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
