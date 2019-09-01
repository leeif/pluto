package migrate

import (
	"github.com/jinzhu/gorm"
	"github.com/leeif/pluto/models"
)

var migrations = []Migrations{
	{
		name:     "change_user_mail",
		function: changeUserMail,
	},
	{
		name:     "change_user_password",
		function: changeUserPassword,
	},
}

func changeUserMail(db *gorm.DB, name string) error {
	// default null
	err := db.Model(&models.User{}).ModifyColumn("mail", "varchar(255)").Error
	if err != nil {
		return err
	}

	// remove unique index
	db.Model(&models.User{}).RemoveIndex("mail")
	return nil
}

func changeUserPassword(db *gorm.DB, name string) error {
	// default null
	err := db.Model(&models.User{}).ModifyColumn("password", "varchar(255)").Error
	if err != nil {
		return err
	}
	return nil
}
