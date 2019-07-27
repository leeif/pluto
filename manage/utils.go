package manage

import (
	"github.com/jinzhu/gorm"
)

const (
	ServerError = iota
	ReqError
)

type PlutoError struct {
	Type int
	Err  error
}

func newPlutoError(t int, err error) *PlutoError {
	return &PlutoError{
		Type: t,
		Err:  err,
	}
}

func create(db *gorm.DB, record interface{}) *PlutoError {
	if err := db.Create(record).Error; err != nil {
		return newPlutoError(ServerError, err)
	}
	return nil
}

func update(db *gorm.DB, record interface{}) *PlutoError {
	if err := db.Save(record).Error; err != nil {
		return newPlutoError(ServerError, err)
	}
	return nil
}
