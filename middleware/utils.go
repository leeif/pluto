package middleware

import (
	"net/http"

	"github.com/urfave/negroni"
)

func TokenVerifyMiddleware(handlers ...negroni.HandlerFunc) http.Handler {
	ng := negroni.New()
	for _, handler := range handlers {
		ng.UseFunc(handler)
	}
	return ng
}

func SessionVerifyMiddleware(handlers ...negroni.HandlerFunc) http.Handler {
	ng := negroni.New()
	for _, handler := range handlers {
		ng.UseFunc(handler)
	}
	return ng
}

func NoVerifyMiddleware(handlers ...negroni.HandlerFunc) http.Handler {
	ng := negroni.New()
	for _, handler := range handlers {
		ng.UseFunc(handler)
	}

	return ng
}
