package models

import (
	"github.com/jinzhu/gorm"
)

type RefreshToken struct {
	gorm.Model
	RefreshToken string `gorm:"type:varchar(255);size:255"`
}
