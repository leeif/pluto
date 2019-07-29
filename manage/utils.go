package manage

import (
	"github.com/jinzhu/gorm"
	"github.com/leeif/pluto/datatype"
)

func create(db *gorm.DB, record interface{}) *datatype.PlutoError {
	if err := db.Create(record).Error; err != nil {
		return datatype.NewPlutoError(datatype.ServerError, err)
	}
	return nil
}

func update(db *gorm.DB, record interface{}) *datatype.PlutoError {
	if err := db.Save(record).Error; err != nil {
		return datatype.NewPlutoError(datatype.ServerError, err)
	}
	return nil
}
