package v1

import (
	"net/http"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/datatype/request"
	routeUtils "github.com/leeif/pluto/utils/route"
	"github.com/leeif/pluto/utils/rsa"
)

func (router *Router) RefreshToken(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	rat := request.RefreshAccessToken{}
	if err := routeUtils.GetRequestData(r, &rat); err != nil {
		return err
	}

	res, err := router.manager.RefreshAccessToken(rat)

	if err != nil {
		return err
	}

	if err := routeUtils.ResponseOK(res, w); err != nil {
		return err
	}

	return nil
}

func (router *Router) PublicKey(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	res := make(map[string]string)
	pbkey, err := rsa.GetPublicKey()

	if err != nil {
		return perror.ServerError.Wrapper(err)
	}

	res["public_key"] = pbkey

	if err := routeUtils.ResponseOK(res, w); err != nil {
		return err
	}

	return nil
}

func (router *Router) VerifyAccessToken(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	accessToken := &request.VerifyAccessToken{}
	if err := routeUtils.GetRequestData(r, accessToken); err != nil {
		return err
	}

	res, err := router.manager.VerifyAccessToken(accessToken.Token)

	if err != nil {
		return err
	}

	if err := routeUtils.ResponseOK(res, w); err != nil {
		return err
	}

	return nil
}
