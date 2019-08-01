package middleware

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/urfave/negroni"
)

type Middleware struct {
	Logger log.Logger
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
