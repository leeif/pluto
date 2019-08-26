package middleware

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/urfave/negroni"

	perror "github.com/leeif/pluto/datatype/pluto_error"

	"github.com/gorilla/context"
)

func NewLogger(logger log.Logger) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		plutoError := context.Get(r, "pluto_error")
		// if plutoError == nil {
		// 	return
		// }
		pe := plutoError.(*perror.PlutoError)
		if pe.LogError != nil {
			level.Error(logger).Log("msg", pe.LogError.Error())
		}
		if pe.HTTPError != nil {
			level.Debug(logger).Log("msg", pe.HTTPError.Error())
		}
	}
}
