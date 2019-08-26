package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-kit/kit/log/level"

	"github.com/go-kit/kit/log"
	"github.com/leeif/pluto/config"
	plog "github.com/leeif/pluto/log"
	"github.com/leeif/pluto/route"
	"github.com/urfave/negroni"
)

type Server struct {
}

func (s Server) RunServer() error {
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
			return err
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

	//start server background
	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			// Error starting or closing listener:
			level.Error(logger).Log("msg", "Server closed with error:"+err.Error())
			os.Exit(1)
		}
	}()

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, os.Interrupt)
	level.Error(logger).Log("msg", "SIGNAL "+(<-quit).String()+" received, then shutting down...")

	// timeout 60s
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		// Error from closing listeners, or context timeout:
		level.Error(logger).Log("msg", "Failed to gracefully shutdown:"+err.Error())
	}
	level.Error(logger).Log("msg", "Server shutdown")
	return nil
}
