package route

import (
	"net/http"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/utils/rsa"
)

func (router *Router) refreshToken(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	rat := request.RefreshAccessToken{}
	if err := getBody(r, &rat); err != nil {
		return err
	}

	res, err := router.manager.RefreshAccessToken(rat)

	if err != nil {
		return err
	}

	if err := responseOK(res, w); err != nil {
		router.logger.Error(err.Error())
	}

	return nil
}

func (router *Router) publicKey(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	res := make(map[string]string)
	pbkey, err := rsa.GetPublicKey()

	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	res["public_key"] = pbkey

	if err := responseOK(res, w); err != nil {
		router.logger.Error(err.Error())
	}

	return nil
}
