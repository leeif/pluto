package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/fx"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/log"
)

type Server struct {
}

func NewMux(lc fx.Lifecycle, config *config.Config, logger *log.PlutoLog) *mux.Router {
	address := ":" + config.Server.Port.String()
	router := mux.NewRouter()
	srv := &http.Server{
		Addr:    address,
		Handler: router,
	}
	lc.Append(fx.Hook{
		// To mitigate the impact of deadlocks in application startup and
		// shutdown, Fx imposes a time limit on OnStart and OnStop hooks. By
		// default, hooks have a total of 30 seconds to complete. Timeouts are
		// passed via Go's usual context.Context.
		OnStart: func(context.Context) error {
			logger.Info("Starting Pluto server at " + address)
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Pluto server")
			return srv.Shutdown(ctx)
		},
	})
	return router
}
