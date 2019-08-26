package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Mail     *string `gorm:"type:varchar(255);size(60);unique;not null"`
	Name     *string `gorm:"type:varchar(60);size(60);not null"`
	Role     string  `gorm:"type:varchar(60);size(60)" json:"-"`
	Gender   *string `gorm:"type:varchar(10);size(10);"`
	Password *string `gorm:"type:varchar(255);not null" json:"-"`
	Birthday *time.Time
	Avatar   string `gorm:"type:varchar(255)"`
	Verified bool   `json:"-"`
}
