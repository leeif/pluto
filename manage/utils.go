package manage

import (
	"github.com/jinzhu/gorm"
	"github.com/leeif/pluto/config"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/log"
)

type Manger struct {
	logger *log.PlutoLog
	config *config.Config
	db     *gorm.DB
}

func NewManager(db *gorm.DB, config *config.Config, logger *log.PlutoLog) *Manger {
	return &Manger{
		logger: logger.With("compoment", "manager"),
		config: config,
		db:     db,
	}
}

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
