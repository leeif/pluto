package middleware

import (
	"net/http"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/log"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

type Middleware struct {
	Logger *log.PlutoLog
	Config *config.Config
}

func (middleware *Middleware) AccessTokenAuthMiddleware(handlers ...negroni.HandlerFunc) http.Handler {
	ng := negroni.New()
	// cors
	c := cors.New(cors.Options{
		AllowedOrigins: middleware.Config.Cros.AllowedOrigins,
	})
	ng.Use(c)
	ng.UseFunc(AccessTokenAuth)
	for _, handler := range handlers {
		ng.UseFunc(handler)
	}
	ng.UseFunc(Logger(middleware.Logger))
	return ng
}

func (middleware *Middleware) AdminAuthMiddleware(handlers ...negroni.HandlerFunc) http.Handler {
	ng := negroni.New()
	// cors
	c := cors.New(cors.Options{
		AllowedOrigins: middleware.Config.Cros.AllowedOrigins,
	})
	ng.Use(c)
	ng.UseFunc(PlutoAdmin())
	for _, handler := range handlers {
		ng.UseFunc(handler)
	}
	ng.UseFunc(Logger(middleware.Logger))
	return ng
}

func (middleware *Middleware) NoAuthMiddleware(handlers ...negroni.HandlerFunc) http.Handler {
	ng := negroni.New()
	// cors
	c := cors.New(cors.Options{
		AllowedOrigins: middleware.Config.Cros.AllowedOrigins,
	})
	ng.Use(c)
	for _, handler := range handlers {
		ng.UseFunc(handler)
	}
	ng.UseFunc(Logger(middleware.Logger))
	return ng
}

func NewMiddle(logger *log.PlutoLog, config *config.Config) *Middleware {
	return &Middleware{
		Logger: logger.With("componment", "middleware"),
		Config: config,
	}
}
