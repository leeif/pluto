package route

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/leeif/pluto/database"
	"github.com/leeif/pluto/manage"
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

		if err := manage.VerifyMail(db, token); err != nil {
			// set err to context for log
			context.Set(r, "pluto_error", err)
			data.Successed = false
		} else {
			data.Successed = true
		}
		if err := responseHTML("register_verify.html", data, w); err != nil {
			fmt.Println(err)
		}
		next(w, r)
	})).Methods("GET")
}
