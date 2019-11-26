package manage

import (
	"database/sql"
	"math/rand"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/log"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func randomToken(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

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
