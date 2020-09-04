package v1

import (
	"github.com/MuShare/pluto/config"
	"github.com/MuShare/pluto/log"
	"github.com/MuShare/pluto/manage"
)

type Router struct {
	manager *manage.Manager
	config  *config.Config
	logger  *log.PlutoLog
}

func NewRouter(manager *manage.Manager, config *config.Config, logger *log.PlutoLog) *Router {
	return &Router{
		manager: manager,
		config:  config,
		logger:  logger,
	}
}
