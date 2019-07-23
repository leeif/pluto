package manage

import (
	"github.com/leeif/pluto/models"
)

func Login() error {
	db, err := models.GetDatabase()
	if err != nil {
		return err
	}
	return nil
}

func Register() error {
	return nil
}

func ResetPassword() error {
	return nil
}