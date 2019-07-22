package models

import (
	"github.com/jinzhu/gorm"
)

func NewDatabase() error {
	// config := config.GetConfig()
	db, err := gorm.Open("")
	if err != nil {
		return err
	}
	db.DB().SetMaxIdleConns(10)
	return nil
}
