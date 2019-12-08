package route

import (
	"net/http"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/utils/rsa"

	"github.com/gorilla/context"
)

func (router *Router) refreshToken(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	rat := request.RefreshAccessToken{}
	if err := getBody(r, &rat); err != nil {
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	res, err := router.manager.RefreshAccessToken(rat)

	if err != nil {
		// set err to context for log
		context.Set(r, "pluto_error", err)
		responseError(err, w)
		next(w, r)
		return
	}

	responseOK(res, w)
}

func (router *Router) publicKey(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	res := make(map[string]string)
	pbkey, err := rsa.GetPublicKey()

	if err != nil {
		perr := perror.ServerError.Wrapper(err)
		responseError(perr, w)
		next(w, r)
		return
	}

	res["public_key"] = pbkey
	responseOK(res, w)
}
