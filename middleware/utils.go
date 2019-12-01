package middleware

import (
	"net/http"

	"github.com/leeif/pluto/log"
	"github.com/urfave/negroni"
)

type Middleware struct {
	Logger *log.PlutoLog
}

func (middleware *Middleware) AccessTokenAuthMiddleware(handlers ...negroni.HandlerFunc) http.Handler {
	ng := negroni.New()
	ng.UseFunc(AccessTokenAuth)
	for _, handler := range handlers {
		ng.UseFunc(handler)
	}
	ng.UseFunc(Logger(middleware.Logger))
	return ng
}

func (middleware *Middleware) AdminAuthMiddleware(handlers ...negroni.HandlerFunc) http.Handler {
	ng := negroni.New()
	ng.UseFunc(PlutoAdmin())
	for _, handler := range handlers {
		ng.UseFunc(handler)
	}
	ng.UseFunc(Logger(middleware.Logger))
	return ng
}

func (middleware *Middleware) NoAuthMiddleware(handlers ...negroni.HandlerFunc) http.Handler {
	ng := negroni.New()
	for _, handler := range handlers {
		ng.UseFunc(handler)
	}
	ng.UseFunc(Logger(middleware.Logger))
	return ng
}

func NewMiddle(logger *log.PlutoLog) *Middleware {
	return &Middleware{
		Logger: logger.With("componment", "middleware"),
	}
}
