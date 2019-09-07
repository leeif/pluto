package config_test

import (
	"os"
	"testing"

	"github.com/leeif/pluto/config"
	"github.com/stretchr/testify/assert"
)

func writeConfigFile(path string, s string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write([]byte(s))
	if err != nil {
		return err
	}
	return nil
}

func deleteConfigFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func TestConfigDefault(t *testing.T) {
	args := []string{"test", "--mail.smtp", "test.smtp.com"}
	// Init config file from command line and config file
	c, err := config.NewConfig(args, "")

	if err != nil {
		t.Fatalf("Expect no err, but: %v", err)
	}

	assert.Equal(t, *c.ConfigFile, "./config.json", "default of config file should be ./config.json")

	assert.Equal(t, "mysql", c.Database.Type.String(), "default of database type should be mysql")
	assert.Equal(t, "127.0.0.1", *c.Database.Host, "default of database host should be 127.0.0.1")
	assert.Equal(t, "root", *c.Database.User, "default of database user should be root")
	assert.Equal(t, "3306", c.Database.Port.String(), "default of database port should be 3306")
	assert.Equal(t, "", *c.Database.Password, "default of database port should be empty")
	assert.Equal(t, "", *c.Database.DB, "default of database port should be empty")
	assert.Equal(t, "test.smtp.com", c.Mail.SMTP.String(), "default of database port should be empty")
}

func TestConfigCustom(t *testing.T) {

	if err := writeConfigFile("./config_test.json", `{
		"server": {
			"port": "8080"
		},
		"log": {
			"file": "/tmp/pluto-test.log",
			"level": "error"
		},
		"rsa": {
			"name": "rsa_test",
			"path": "/tmp"
		},
		"database": {
			"type": "mysql",
			"host": "192.168.1.1",
			"user": "www",
			"port": "3306",
			"password": "www",
			"db": "pluto_server"
		},
		"mail": {
			"smtp": "test.smtp.com"
		}
	}`); err != nil {
		t.Fatal(err)
		return
	}

	defer deleteConfigFile("./config_test.json")

	args := []string{"test", "--config.file", "./config_test.json"}
	// Init config file from command line and config file
	c, err := config.NewConfig(args, "")
	if err != nil {
		t.Fatalf("Expect no error, but: %v", err)
	}

	assert.Equal(t, *c.ConfigFile, "./config_test.json", "config file should be ./config.json")

	assert.Equal(t, c.Server.Port.String(), "8080", "server port should be 8080")

	assert.Equal(t, c.Log.Level.String(), "error", "log level should be error")
	assert.Equal(t, *c.Log.File, "/tmp/pluto-test.log", "server port should be 8080")

	assert.Equal(t, *c.RSA.Name, "rsa_test", "rsa file name should be rsa_test")
	assert.Equal(t, *c.RSA.Path, "/tmp", "server port should be /tmp")

	assert.Equal(t, "mysql", c.Database.Type.String(), "database type should be mysql")
	assert.Equal(t, "192.168.1.1", *c.Database.Host, "database host should be 127.0.0.1")
	assert.Equal(t, "www", *c.Database.User, "database user should be root")
	assert.Equal(t, "3306", c.Database.Port.String(), "database port should be 3306")
	assert.Equal(t, "www", *c.Database.Password, "database port should be empty")
	assert.Equal(t, "pluto_server", *c.Database.DB, "database port should be empty")

	assert.Equal(t, "test.smtp.com", c.Mail.SMTP.String(), "database port should be empty")
}
