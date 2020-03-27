package middleware

import (
	"encoding/json"
	"github.com/urfave/negroni"
	"net/http"
	"time"

	"github.com/gorilla/context"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/utils/jwt"
)

func AccessTokenAuth() negroni.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		accessToken, perr := getAuthorizationHeader(r)
		if perr != nil {
			context.Set(r, "pluto_error", perr)
			next(w, r)
			return
		}

		jwtToken, perr := jwt.VerifyJWT(accessToken)
		if perr != nil {
			context.Set(r, "pluto_error", perr)
			next(w, r)
			return
		}

		accessPayload := &jwt.AccessPayload{}

		if err := json.Unmarshal(jwtToken.Payload, &accessPayload); err != nil {
			context.Set(r, "pluto_error", perror.ServerError.Wrapper(err))
			next(w, r)
			return
		}

		if accessPayload.Type != jwt.ACCESS {
			context.Set(r, "pluto_error", perror.InvalidJWTToken)
			next(w, r)
			return
		}

		if time.Now().Unix() > accessPayload.Expire {
			context.Set(r, "pluto_error", perror.JWTTokenExpired)
			next(w, r)
			return
		}

		context.Set(r, "payload", accessPayload)

		next(w, r)
	}
}
