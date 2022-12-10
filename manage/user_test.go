package manage_test

import (
	"testing"

	"github.com/MuShare/pluto/config"
	"github.com/MuShare/pluto/database"
	"github.com/MuShare/pluto/datatype/request"
	plog "github.com/MuShare/pluto/log"
	"github.com/MuShare/pluto/manage"
	"github.com/MuShare/pluto/utils/jwt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func getManager() (*manage.Manager, error) {
	args := []string{"test", "--database.password", "12345678"}
	c, err := config.NewConfig(args, "")
	if err != nil {
		return nil, err
	}
	d, err := database.NewDatabase(c)
	if err != nil {
		return nil, err
	}
	pl, err := plog.NewLogger(c)
	if err != nil {
		return nil, err
	}

	return manage.NewManager(d, c, pl)
}

func TestRandomUserName(t *testing.T) {
	m, err := getManager()
	if err != nil {
		t.Fatal(err)
	}
	if m == nil {
		t.Fatal("manager is nil")
	}

	userName, err2 := m.RandomUserName("hello")
	if err2 != nil {
		t.Fatal(err2)
	}

	assert.Equal(t, len("hello")+6, len(userName), " Random user name length should be equal to length of prefix + 6")
}

func TestCreateApplication(t *testing.T) {
	m, err := getManager()
	if err != nil {
		t.Fatal(err)
	}
	if m == nil {
		t.Fatal("manager is nil")
	}
	request := request.CreateApplication{
		Name: "easy-pluto",
	}
	app, err2 := m.CreateApplication(request)
	if err2 != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "easy-pluto", app.Name, "app name should be easy-pluto")
}

func TestRegisterWithEmail(t *testing.T) {
	m, err := getManager()
	if err != nil {
		t.Fatal(err)
	}
	if m == nil {
		t.Fatal("manager is nil")
	}
	request := request.MailRegister{
		Mail:     "yanyin1986@gmail.com",
		UserID:   "yanyin1986",
		Name:     "yanyin1986",
		Password: "12345678",
		AppName:  "easy-pluto",
	}
	user, err2 := m.RegisterWithEmail(request, false)
	if err2 != nil {
		t.Fatal(err2)
	}
	assert.Equal(t, "yanyin1986", user.Name, "user name should be yanyin1986")
}

func TestDeleteUser(t *testing.T) {
	m, err := getManager()
	if err != nil {
		t.Fatal(err)
	}
	if m == nil {
		t.Fatal("manager is nil")
	}
	/*
		request := request.DeleteUser{
			UserID: "yanyin1986",
		}
	*/
	payload := jwt.AccessPayload{
		UserID: 1,
		AppID:  "easy-pluto",
		Scopes: []string{"user:read"},
	}
	err2 := m.DeleteUser(&payload)
	if err2 != nil {
		t.Fatal(err2)
	}
	assert.Equal(t, true, "user should be deleted")
}
