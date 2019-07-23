package models

import (
	"github.com/jinzhu/gorm"
)

type Device struct {
	gorm.Model
	Identifier string `gorm:"type:varchar;unique;not null"`
	Agent      string `gorm:"type:varchar;not null"`
}
