package route

import (
	"net/http"
)

func (router *Router) healthCheck(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	respBody := make(map[string]interface{})
	respBody["version"] = router.config.Version
	responseOK(respBody, w)
}
