package models

import (
	"github.com/jinzhu/gorm"
)

// Application : applications registered in auth server
type Application struct {
	gorm.DB
	Name     string `gorm:"type:varchar(100);size:100;not null"`
	Callback string `gorm:"type:varchar(255);size:255;not null"`
}
