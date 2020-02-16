package server

import (
	"context"
	"github.com/micro/go-micro/server"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.uber.org/fx"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/log"
	"github.com/micro/go-plugins/registry/consul"
	httpServer "github.com/micro/go-plugins/server/http"

	"github.com/micro/go-micro"
)

type Server struct {
}

func NewMux(lc fx.Lifecycle, config *config.Config, logger *log.PlutoLog) *mux.Router {
	address := ":" + config.Server.Port.String()
	router := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins: config.Cors.AllowedOrigins,
		AllowedHeaders: config.Cors.AllowedHeaders,
	})
	handler := c.Handler(router)

	// microservices
	registry := consul.NewRegistry()
	srv := httpServer.NewServer(
		server.Name("pluto"),
		server.Address(address),
	)
	hd := srv.NewHandler(handler)
	srv.Handle(hd)
	ctx, cancel := context.WithCancel(context.Background())
	service := micro.NewService(
		micro.Server(srv),
		micro.Name("pluto"),
		micro.Registry(registry),
		micro.Context(ctx),
	)
	lc.Append(fx.Hook{
		// To mitigate the impact of deadlocks in application startup and
		// shutdown, Fx imposes a time limit on OnStart and OnStop hooks. By
		// default, hooks have a total of 30 seconds to complete. Timeouts are
		// passed via Go's usual context.Context.
		OnStart: func(ctx context.Context) error {
			logger.Info("Starting Pluto server at " + address)
			go service.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("Stopping Pluto server")
			cancel()
			// wait for consul deregister
			time.Sleep(5*time.Second)
			return nil
		},
	})
	return router
}
