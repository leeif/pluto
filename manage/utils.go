package manage

import (
	"database/sql"

	"github.com/leeif/pluto/config"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/models"
	"github.com/volatiletech/sqlboiler/boil"
)

type Manager struct {
	logger *log.PlutoLog
	config *config.Config
	db     *sql.DB
}

func NewManager(db *sql.DB, config *config.Config, logger *log.PlutoLog) *Manager {
	return &Manager{
		logger: logger.With("compoment", "manager"),
		config: config,
		db:     db,
	}
}

const (
	OperationMailLogin   = "maillogin"
	OperationGoogleLogin = "googlelogin"
	OperationWechatLogin = "wechatlogin"
	OperationAppleLogin  = "applelogin"

	OperationLogout        = "logout"
	OperationResetPassword = "reset_password"
	OperationRefreshToken  = "refresh_token"
)

func historyOperation(tx *sql.Tx, operationType string, userID uint) *perror.PlutoError {
	historyOperation := models.HistoryOperation{}
	historyOperation.UserID = userID
	historyOperation.Type = operationType
	if err := historyOperation.Insert(tx, boil.Infer()); err != nil {
		return perror.ServerError.Wrapper(err)
	}
	return nil
}
