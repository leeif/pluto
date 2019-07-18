package models

import (
	"github.com/jinzhu/gorm"
)

type Migration struct {
	gorm.Model
	Name string `gorm:"type:varchar(100);size:100"`
}
