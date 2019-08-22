package route

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/leeif/pluto/database"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/manage"
	"github.com/leeif/pluto/utils/jwt"
)

func (route *Route) webRoute(router *mux.Router) {
	db, err := database.GetDatabase()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	router.Handle("/mail/verify/{token}", route.middleware.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		vars := mux.Vars(r)
		token := vars["token"]

		type Data struct {
			Successed bool
		}
		data := &Data{}

		if err := manage.RegisterVerify(db, token); err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			data.Successed = false
		} else {
			data.Successed = true
		}

		responseHTML("register_verify_result.html", data, w)
		next(w, r)
	})).Methods("GET")

	router.Handle("/password/reset/{token}", route.middleware.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		vars := mux.Vars(r)
		token := vars["token"]
		if err := manage.ResetPasswordPage(db, token); err != nil {
			context.Set(r, "pluto_error", err)
			responseHTML("error/oops.html", nil, w)
			next(w, r)
			return
		}

		type Data struct {
			Token string
		}
		data := &Data{Token: token}

		responseHTML("password_reset.html", data, w)
		next(w, r)
	})).Methods("GET")

	router.Handle("/password/reset/result/{token}", route.middleware.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		vars := mux.Vars(r)
		token := vars["token"]

		header, payload, err := jwt.VerifyB64JWT(token)
		// token verify failed
		if err != nil {
			context.Set(r, "pluto_error", err)
			responseHTML("error/oops.html", nil, w)
			next(w, r)
			return
		}

		head := jwt.Head{}
		json.Unmarshal(header, &head)
		if head.Type != jwt.PASSWORDRESETRESULT {
			context.Set(r, "pluto_error", perror.InvalidJWTToekn)
			responseHTML("error/oops.html", nil, w)
			next(w, r)
			return
		}

		prp := jwt.PasswordResetResultPayload{}
		json.Unmarshal(payload, &prp)

		type Data struct {
			Message string
		}
		data := &Data{Message: prp.Message}

		responseHTML("password_reset_result.html", data, w)
		next(w, r)
	})).Methods("GET")
}
