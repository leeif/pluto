package models

import (
	"github.com/jinzhu/gorm"
)

type Salt struct {
	gorm.Model
	UserID uint   `gorm:"column:user_id"`
	Salt   string `gorm:"type:varchar(255);size:255;not null"`
}
