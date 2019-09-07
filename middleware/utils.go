package middleware

import (
	"net/http"

	"github.com/leeif/pluto/log"
	"github.com/urfave/negroni"
)

type Middleware struct {
	Logger *log.PlutoLog
}

func (middleware *Middleware) TokenVerifyMiddleware(handlers ...negroni.HandlerFunc) http.Handler {
	ng := negroni.New()
	for _, handler := range handlers {
		ng.UseFunc(handler)
	}
	ng.UseFunc(NewLogger(middleware.Logger))
	return ng
}

func (middleware *Middleware) SessionVerifyMiddleware(handlers ...negroni.HandlerFunc) http.Handler {
	ng := negroni.New()
	for _, handler := range handlers {
		ng.UseFunc(handler)
	}
	return ng
}

func (middleware *Middleware) NoVerifyMiddleware(handlers ...negroni.HandlerFunc) http.Handler {
	ng := negroni.New()
	for _, handler := range handlers {
		ng.UseFunc(handler)
	}
	ng.UseFunc(NewLogger(middleware.Logger))
	return ng
}

func NewMiddle(logger *log.PlutoLog) *Middleware {
	return &Middleware{
		Logger: logger.With("componment", "middleware"),
	}
}
