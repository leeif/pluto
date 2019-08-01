package middleware

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/urfave/negroni"

	"github.com/leeif/pluto/datatype"

	"github.com/gorilla/context"
)

func NewLogger(logger log.Logger) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		plutoError := context.Get(r, "pluto_err")
		if plutoError == nil {
			return
		}
		pe := plutoError.(datatype.PlutoError)
		switch pe.Type {
		case datatype.ServerError:
			level.Error(logger).Log("msg", pe.Err.Error())
		case datatype.ReqError:
			level.Debug(logger).Log("msg", pe.Err.Error())
		}
	}
}
