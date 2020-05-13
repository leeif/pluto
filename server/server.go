package server

import (
	"context"
	"net/http"
	"time"

	"log"

	"github.com/gorilla/mux"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/registry/mdns"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-plugins/registry/consul"
	"github.com/rs/cors"
	"go.uber.org/fx"

	"github.com/leeif/pluto/config"
	httpServer "github.com/micro/go-plugins/server/http"
)

type Server struct {
}

func NewMux(lc fx.Lifecycle, config *config.Config) *mux.Router {
	address := ":" + config.Server.Port.String()
	router := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedOrigins: config.Cors.AllowedOrigins,
		AllowedHeaders: config.Cors.AllowedHeaders,
	})

	// microservices
	var rgy registry.Registry
	if config.Registry.Consul {
		rgy = consul.NewRegistry(
			registry.Addrs(config.Registry.ConsulAddress + ":" + config.Registry.ConsulPort.String()),
		)
	} else {
		rgy = mdns.NewRegistry()
	}

	srv := httpServer.NewServer(
		server.Name(config.Server.ServerName),
		server.Address(address),
	)

	handler := c.Handler(router)
	hd := srv.NewHandler(handler)
	srv.Handle(hd)
	ctx, cancel := context.WithCancel(context.Background())
	service := micro.NewService(
		micro.Server(srv),
		micro.Name(config.Registry.ServiceName),
		micro.Context(ctx),
		micro.Registry(rgy),
	)

	lc.Append(fx.Hook{
		// To mitigate the impact of deadlocks in application startup and
		// shutdown, Fx imposes a time limit on OnStart and OnStop hooks. By
		// default, hooks have a total of 30 seconds to complete. Timeouts are
		// passed via Go's usual context.Context.
		OnStart: func(context.Context) error {
			log.Println("Starting Pluto server at " + address)
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go service.Run()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping Pluto server")
			cancel()
			// wait for consul deregister
			time.Sleep(5 * time.Second)
			return nil
		},
	})
	return router
}
