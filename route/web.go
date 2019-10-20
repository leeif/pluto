package route

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/leeif/pluto/config"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/manage"
	"github.com/leeif/pluto/middleware"
	"github.com/leeif/pluto/utils/jwt"
)

func webRouter(router *mux.Router, db *gorm.DB, config *config.Config, logger *log.PlutoLog) {

	mw := middleware.NewMiddle(logger)
	manager := manage.NewManager(db, config, logger)

	router.Handle("/mail/verify/{token}", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		vars := mux.Vars(r)
		token := vars["token"]

		type Data struct {
			Successed bool
		}
		data := &Data{}

		if err := manager.RegisterVerify(token); err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			next(w, r)
			data.Successed = false
			responseHTMLError("register_verify_result.html", data, w, http.StatusForbidden)
		} else {
			data.Successed = true
			responseHTMLOK("register_verify_result.html", data, w)
		}

	})).Methods("GET")

	router.Handle("/password/reset/{token}", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		vars := mux.Vars(r)
		token := vars["token"]
		if err := manager.ResetPasswordPage(token); err != nil {
			context.Set(r, "pluto_error", err)
			next(w, r)
			responseHTMLError("error.html", nil, w, http.StatusInternalServerError)
			return
		}

		type Data struct {
			Token string
		}
		data := &Data{Token: token}

		responseHTMLOK("password_reset.html", data, w)
	})).Methods("GET")

	router.Handle("/password/reset/result/{token}", mw.NoVerifyMiddleware(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		vars := mux.Vars(r)
		token := vars["token"]

		jwtToken, err := jwt.VerifyB64JWT(token)
		// token verify failed
		if err != nil {
			context.Set(r, "pluto_error", err)
			next(w, r)
			responseHTMLError("error.html", nil, w, http.StatusInternalServerError)
			return
		}

		head := jwt.Head{}
		json.Unmarshal(jwtToken.Head, &head)
		if head.Type != jwt.PASSWORDRESETRESULT {
			context.Set(r, "pluto_error", perror.InvalidJWTToekn)
			next(w, r)
			responseHTMLError("error.html", nil, w, http.StatusInternalServerError)
			return
		}

		prp := jwt.PasswordResetResultPayload{}
		json.Unmarshal(jwtToken.Payload, &prp)

		type Data struct {
			Successed bool
		}
		data := &Data{Successed: prp.Successed}
		responseHTMLOK("password_reset_result.html", data, w)
	})).Methods("GET")
}
