package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Mail     string `gorm:"type:varchar(254);unique;not null"`
	Name     string `gorm:"type:varchar(60);not null"`
	Role     string `gorm:"type:varchar(60);not null"`
	Password string `gorm:"type:varchar(60);not null"`
}
