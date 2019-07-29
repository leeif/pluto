package server

import (
	"net/http"

	"github.com/go-kit/kit/log/level"

	"github.com/go-kit/kit/log"
	"github.com/gorilla/mux"
	"github.com/leeif/pluto/config"
	"github.com/urfave/negroni"
)

type Server struct {
	config *config.ServerConfig
	logger log.Logger
	router *mux.Router
}

func (s Server) RunServer() error {
	address := ":" + s.config.Port.String()

	n := negroni.New()

	n.UseHandler(s.router)

	level.Info(s.logger).Log("msg", "Start pluto server at "+address)
	err := http.ListenAndServe(address, n)
	if err != nil {
		return err
	}

	return nil
}

func NewServer(router *mux.Router, logger log.Logger, config *config.ServerConfig) *Server {
	l := log.With(logger, "component", "server")

	server := Server{
		config: config,
		logger: l,
		router: router,
	}
	return &server
}
