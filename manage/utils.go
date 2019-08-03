package manage

import (
	"github.com/jinzhu/gorm"
	perror "github.com/leeif/pluto/datatype/pluto_error"
)

func create(db *gorm.DB, record interface{}) *perror.PlutoError {
	if err := db.Create(record).Error; err != nil {
		return perror.NewServerError(err)
	}
	return nil
}

func update(db *gorm.DB, record interface{}) *perror.PlutoError {
	if err := db.Save(record).Error; err != nil {
		return perror.NewServerError(err)
	}
	return nil
}
