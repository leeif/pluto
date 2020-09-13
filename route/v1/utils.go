package v1

import (
	"github.com/MuShare/pluto/config"
	"github.com/MuShare/pluto/log"
	"github.com/MuShare/pluto/manage"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type Router struct {
	manager *manage.Manager
	config  *config.Config
	logger  *log.PlutoLog
	bundle	*i18n.Bundle
}

func NewRouter(manager *manage.Manager, config *config.Config, logger *log.PlutoLog, bundle *i18n.Bundle) *Router {
	return &Router{
		manager: manager,
		config:  config,
		logger:  logger,
		bundle:  bundle,
	}
}
