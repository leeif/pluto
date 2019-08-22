package server

import (
	"net/http"
	"os"

	"github.com/go-kit/kit/log/level"

	"github.com/go-kit/kit/log"
	"github.com/leeif/pluto/config"
	plog "github.com/leeif/pluto/log"
	"github.com/leeif/pluto/route"
	"github.com/urfave/negroni"
)

type Server struct {
}

func (s Server) RunServer() (*http.Server, error) {
	config := config.GetConfig()
	address := ":" + config.Server.Port.String()

	// set logger
	var logger log.Logger
	var file *os.File
	if *config.Log.File == "" {
		file = os.Stdout
	} else {
		var err error
		file, err = os.OpenFile(*config.Log.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
	}
	defer file.Close()
	logger = plog.GetFileLogger(config.Log, file)

	// get route
	n := negroni.New()
	r := route.Route{
		Logger: logger,
	}
	n.UseHandler(r.GetRouter(logger))

	// start server
	level.Info(logger).Log("msg", "Start pluto server at "+address)
	srv := &http.Server{
		Addr:    address,
		Handler: n,
	}

	return srv, nil
}
