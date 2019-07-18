package server

import (
	"net/http"

	"github.com/go-kit/kit/log/level"

	"github.com/go-kit/kit/log"
	"github.com/julienschmidt/httprouter"
	"github.com/leeif/pluto/api"
	"github.com/leeif/pluto/config"
	plog "github.com/leeif/pluto/log"
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

func NewServer() *Server {
	logger := plog.GetLogger()
	server := Server{
		config: config.GetConfig().Server,
		logger: log.With(logger, "component", "server"),
	}
	return &server
}
