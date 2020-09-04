package middleware

import (
	"net/http"
	"strings"

	perror "github.com/MuShare/pluto/datatype/pluto_error"
	"github.com/urfave/negroni"
)

type HandlerWrapper func(func(http.ResponseWriter, *http.Request) *perror.PlutoError) negroni.HandlerFunc

func AccessTokenAuthMiddleware(handlerWrapper HandlerWrapper, handlers ...func(http.ResponseWriter, *http.Request) *perror.PlutoError) http.Handler {
	ng := negroni.New()
	ng.UseFunc(handlerWrapper(AccessTokenAuth))
	for _, handler := range handlers {
		ng.UseFunc(handlerWrapper(handler))
	}
	return ng
}

func PlutoUserAuthMiddleware(handlerWrapper HandlerWrapper, handlers ...func(http.ResponseWriter, *http.Request) *perror.PlutoError) http.Handler {
	ng := negroni.New()
	ng.UseFunc(handlerWrapper(PlutoUser))
	for _, handler := range handlers {
		ng.UseFunc(handlerWrapper(handler))
	}
	return ng
}

func PlutoAdminAuthMiddleware(handlerWrapper HandlerWrapper, handlers ...func(http.ResponseWriter, *http.Request) *perror.PlutoError) http.Handler {
	ng := negroni.New()
	ng.UseFunc(handlerWrapper(PlutoAdmin))
	for _, handler := range handlers {
		ng.UseFunc(handlerWrapper(handler))
	}
	return ng
}

func NoAuthMiddleware(handlerWrapper HandlerWrapper, handlers ...func(http.ResponseWriter, *http.Request) *perror.PlutoError) http.Handler {
	ng := negroni.New()
	for _, handler := range handlers {
		ng.UseFunc(handlerWrapper(handler))
	}
	return ng
}

func getAccessToken(r *http.Request) (string, *perror.PlutoError) {

	if cookie, err := r.Cookie("access_token"); err == nil {
		jwt := cookie.Value
		return jwt, nil
	} else if err != nil && err != http.ErrNoCookie {
		return "", perror.ServerError.Wrapper(err)
	}

	auth := strings.Fields(r.Header.Get("Authorization"))

	if len(auth) != 2 {
		return "", perror.Unauthorized
	}

	if strings.ToLower(auth[0]) != "jwt" && strings.ToLower(auth[0]) != "bearer" {
		return "", perror.Unauthorized
	}

	return auth[1], nil
}
