package middleware

import (
	"fmt"
	"net/http"

	"github.com/urfave/negroni"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/log"

	"github.com/gorilla/context"
)

func NewLogger(logger *log.PlutoLog) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		plutoError := context.Get(r, "pluto_error")
		if plutoError == nil {
			return
		}
		url := r.URL.String()
		pe := plutoError.(*perror.PlutoError)
		if pe.LogError != nil {
			logger.Error(fmt.Sprintf("[%s]:%s", url, pe.LogError.Error()))
		}
		if pe.HTTPError != nil {
			logger.Debug(fmt.Sprintf("[%s]:%s", url, pe.HTTPError.Error()))
		}
	}
}
