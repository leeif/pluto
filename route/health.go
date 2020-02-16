package route

import (
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"net/http"
)

func (router *Router) healthCheck(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	respBody := make(map[string]interface{})
	respBody["version"] = router.config.Version

	if err := responseOK(respBody, w); err != nil {
		router.logger.Error(err.Error())
	}

	return nil
}
