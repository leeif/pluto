package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/log"

	"github.com/alecthomas/template"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	resp "github.com/leeif/pluto/datatype/response"
)

type router struct {
	Name   string
	Prefix string
	Func   func(router *mux.Router, db *gorm.DB, config *config.Config, logger *log.PlutoLog)
}

var routers = []router{
	{
		Name:   "register",
		Prefix: "/api/user",
		Func:   registerRouter,
	},
	{
		Name:   "login",
		Prefix: "/api/user",
		Func:   loginRouter,
	},
	{
		Name:   "user",
		Prefix: "/api/user",
		Func:   userRouter,
	},
	{
		Name:   "auth",
		Prefix: "/api/auth",
		Func:   authRouter,
	},
	{
		Name:   "web",
		Prefix: "/",
		Func:   webRouter,
	},
	{
		Name:   "healthCheck",
		Prefix: "/",
		Func:   healthCheckRouter,
	},
}

func Router(router *mux.Router, db *gorm.DB, config *config.Config, logger *log.PlutoLog) {

	for _, r := range routers {
		logger.Info(fmt.Sprintf("Register %s router", r.Name))
		sub := router.PathPrefix(r.Prefix).Subrouter()
		r.Func(sub, db, config, logger)
	}
}

func getBody(r *http.Request, revicer interface{}) *perror.PlutoError {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return perror.ServerError.Wrapper(errors.New("Read body failed: " + err.Error()))
	}

	contentType := r.Header.Get("Content-type")
	if strings.Contains(contentType, "application/json") {
		err := json.Unmarshal(body, &revicer)
		if err != nil {
			return perror.BadRequest
		}
	}
	pr, ok := revicer.(request.PlutoRequest)
	// check request body validation
	if ok && !pr.Validation() {
		return perror.BadRequest
	}
	return nil
}

func responseOK(body interface{}, w http.ResponseWriter) error {
	response := resp.ReponseOK{}
	response.Status = resp.STATUSOK
	response.Body = body
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func responseError(plutoError *perror.PlutoError, w http.ResponseWriter) error {
	response := resp.ReponseError{}
	response.Status = resp.STATUSERROR
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(plutoError.HTTPCode)

	m := make(map[string]interface{})
	m["code"] = plutoError.PlutoCode
	m["message"] = plutoError.HTTPError.Error()
	response.Error = m
	b, err := json.Marshal(response)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func responseHTMLOK(file string, data interface{}, w http.ResponseWriter) error {
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	dir, _ := os.Getwd()
	t, err := template.ParseFiles(path.Join(dir, "views", file), path.Join(dir, "views/template", "header.html"))
	if err != nil {
		return err
	}
	err = t.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}

func responseHTMLError(file string, data interface{}, w http.ResponseWriter, status int) error {
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(status)
	dir, _ := os.Getwd()
	t, err := template.ParseFiles(path.Join(dir, "views", file), path.Join(dir, "views/template", "header.html"))
	if err != nil {
		return err
	}
	err = t.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}
