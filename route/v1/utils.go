package v1

import (
	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/manage"
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
