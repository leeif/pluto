package models

import (
	"github.com/jinzhu/gorm"
)

type Salt struct {
	gorm.DB
	salt string `gorm:"type:varchar(255);size:255;not null"`
}

