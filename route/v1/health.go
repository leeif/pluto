package v1

import (
	"net/http"

	perror "github.com/MuShare/pluto/datatype/pluto_error"
	routeUtils "github.com/MuShare/pluto/utils/route"
)

func (router *Router) HealthCheck(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	respBody := make(map[string]interface{})
	respBody["version"] = router.config.Version

	if err := routeUtils.ResponseOK(respBody, w); err != nil {
		return err
	}

	return nil
}
