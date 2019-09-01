package manage

import (
	"github.com/jinzhu/gorm"
	perror "github.com/leeif/pluto/datatype/pluto_error"
)

func create(tx *gorm.DB, record interface{}) *perror.PlutoError {
	if err := tx.Create(record).Error; err != nil {
		tx.Rollback()
		return perror.ServerError.Wrapper(err)
	}
	return nil
}

func update(tx *gorm.DB, record interface{}) *perror.PlutoError {
	if err := tx.Save(record).Error; err != nil {
		tx.Rollback()
		return perror.ServerError.Wrapper(err)
	}
	return nil
}
