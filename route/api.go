package route

import (
	"github.com/leeif/pluto/datatype"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/leeif/pluto/database"

	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/manage"

	"github.com/urfave/negroni"

	"github.com/gorilla/mux"
)

func GetAPIRouter(router *mux.Router) {
	router.PathPrefix("/api").Handler(userRoute())
	router.PathPrefix("/api").Handler(authRoute())
}

func userRoute() http.Handler {
	db, err := database.GetDatabase()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	router := mux.NewRouter()
	router.Handle("/user/register", negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			fmt.Println("register")
		}),
	))
	router.Handle("/user/login", negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			body, _ := ioutil.ReadAll(r.Body)
			contentType := r.Header.Get("Content-type")
			login := request.MailLogin{}
			if contentType == "application/json" {
				json.Unmarshal(body, &login)
			}
			if jwtToken, err := manage.LoginWithEmail(db, login); err != nil {
				response(datatype.STATUSERROR, , w)
			} else {
				m := make(map[string]interface{})
				response(datatype.STATUSOK, , w)
			}
		}),
	))
	return router
}

func authRoute() http.Handler {
	router := mux.NewRouter()
	router.Handle("/user/auth/refresh", negroni.New(
		negroni.HandlerFunc(func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			fmt.Println("refresh")
		}),
	))
	return router
}
