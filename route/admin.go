package route

import (
	"net/http"

	"github.com/gorilla/context"
	"github.com/leeif/pluto/datatype/request"
)

func (router *Router) AdminPage(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

}

func (router *Router) FindUser(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fu := request.FindUser{}
	if err := getQuery(r, &fu); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}
}
