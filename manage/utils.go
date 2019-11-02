package manage

import (
	"github.com/jinzhu/gorm"
	"github.com/leeif/pluto/config"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/models"
)

type Manager struct {
	logger *log.PlutoLog
	config *config.Config
	db     *gorm.DB
}

func NewManager(db *gorm.DB, config *config.Config, logger *log.PlutoLog) *Manager {
	return &Manager{
		logger: logger.With("compoment", "manager"),
		config: config,
		db:     db,
	}
}

func create(tx *gorm.DB, record interface{}) *perror.PlutoError {
	if err := tx.Create(record).Error; err != nil {
		return perror.ServerError.Wrapper(err)
	}
	return nil
}

func update(tx *gorm.DB, record interface{}) *perror.PlutoError {
	if err := tx.Save(record).Error; err != nil {
		return perror.ServerError.Wrapper(err)
	}
	return nil
}

const (
	OperationMailLogin   = "maillogin"
	OperationGoogleLogin = "googlelogin"
	OperationWechatLogin = "wecgatlogin"

	OperationLogout        = "logout"
	OperationResetPassword = "reset_password"
	OperationRefreshToken  = "refresh_token"
)

func historyOperation(tx *gorm.DB, operationType string, userID uint) *perror.PlutoError {
	historyOperation := models.HistoryOperation{
		UserID:        userID,
		OperationType: operationType,
	}
	if err := create(tx, &historyOperation); err != nil {
		return err
	}
	return nil
}
