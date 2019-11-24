package manage

import (
	"database/sql"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/log"
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
