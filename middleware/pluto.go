package middleware

import (
	"net/http"
	"time"

	"github.com/leeif/pluto/utils/general"

	"github.com/gorilla/context"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/utils/jwt"
	"github.com/wxnacy/wgo/arrays"
)

func getAccessPayload(r *http.Request) (*jwt.AccessPayload, *perror.PlutoError) {
	accessToken, perr := getAccessToken(r)

	if perr != nil {
		return nil, perr
	}

	jwtToken, perr := jwt.VerifyRS256JWT(accessToken)
	if perr != nil {
		return nil, perr
	}

	accessPayload := &jwt.AccessPayload{}

	if perr := jwtToken.UnmarshalPayload(accessPayload); perr != nil {
		return nil, perr
	}

	if accessPayload.Type != jwt.ACCESS {
		return nil, perror.InvalidJWTToken
	}

	if time.Now().Unix() > accessPayload.Expire {
		return nil, perror.JWTTokenExpired
	}

	context.Set(r, "payload", accessPayload)

	return accessPayload, nil
}

func PlutoAdmin(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	accessPayload, perr := getAccessPayload(r)

	if perr != nil {
		return nil
	}

	if accessPayload.AppID != general.PlutoApplication {
		return perror.InvalidAccessToken
	}

	if arrays.ContainsString(accessPayload.Scopes, general.PlutoAdminScope) == -1 {
		return perror.NotPlutoAdmin
	}

	return nil
}

func PlutoUser(w http.ResponseWriter, r *http.Request) *perror.PlutoError {
	accessPayload, perr := getAccessPayload(r)

	if perr != nil {
		return nil
	}

	if accessPayload.AppID != general.PlutoApplication {
		return perror.InvalidAccessToken
	}

	if arrays.ContainsString(accessPayload.Scopes, general.PlutoAdminScope) != -1 {
		return nil
	}

	if arrays.ContainsString(accessPayload.Scopes, general.PlutoUserScope) == -1 {
		return perror.NotPlutoUser
	}

	return nil
}
