package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/leeif/pluto/utils/general"

	"github.com/gorilla/context"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/wxnacy/wgo/arrays"
)

func PlutoAdmin(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	accessToken, perr := getAccessToken(r)

	if perr != nil {
		return perr
	}

	jwtToken, perr := jwt.VerifyJWT(accessToken)
	if perr != nil {
		return perr
	}

	accessPayload := &jwt.AccessPayload{}

	if err := json.Unmarshal(jwtToken.Payload, &accessPayload); err != nil {
		return perror.ServerError.Wrapper(err)
	}

	if accessPayload.Type != jwt.ACCESS {
		return perror.InvalidJWTToken
	}

	if time.Now().Unix() > accessPayload.Expire {
		return perror.JWTTokenExpired
	}

	if accessPayload.AppID != general.PlutoAdminApplication {
		return perror.InvalidAccessToken
	}

	if arrays.ContainsString(accessPayload.Scopes, general.PlutoAdminScope) == -1 {
		return perror.NotPlutoAdmin
	}

	context.Set(r, "payload", accessPayload)

	return nil
}
