package middleware

import (
	"net/http"
	"strings"

	"github.com/leeif/pluto/config"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/log"
	"github.com/urfave/negroni"
)

type Middleware struct {
	Logger *log.PlutoLog
	Config *config.Config
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

func NewMiddle(logger *log.PlutoLog, config *config.Config) *Middleware {
	return &Middleware{
		Logger: logger.With("componment", "middleware"),
		Config: config,
	}
}

func getAuthorizationHeader(r *http.Request) (string, *perror.PlutoError) {
	auth := strings.Fields(r.Header.Get("Authorization"))

	if len(auth) != 2 {
		return "", perror.Unauthorized
	}

	if strings.ToLower(auth[0]) != "jwt" && strings.ToLower(auth[0]) != "bearer" {
		return "", perror.Unauthorized
	}

	return auth[1], nil
}
