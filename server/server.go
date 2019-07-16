package server

import (
	"net/http"

	"github.com/go-kit/kit/log/level"

	"github.com/go-kit/kit/log"
	"github.com/julienschmidt/httprouter"
	"github.com/leeif/pluto/api"
	"github.com/leeif/pluto/config"
	"github.com/urfave/negroni"
)

type Server struct {
	config *config.ServerConfig
	logger log.Logger
}

func (s Server) RunServer() error {
	address := ":" + s.config.Port.String()

	router := httprouter.New()
	api.AddRoute(router)

	n := negroni.Classic()
	n.Use(negroni.NewLogger())

	n.UseHandler(router)

	level.Info(s.logger).Log("msg", "Start pluto server at "+address)
	err := http.ListenAndServe(address, n)
	if err != nil {
		return err
	}

	return nil
}

func NewServer(config *config.ServerConfig, logger log.Logger) *Server {
	server := Server{
		config: config,
		logger: log.With(logger, "component", "server"),
	}
	return &server
}
