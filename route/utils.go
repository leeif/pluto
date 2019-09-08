package route

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/leeif/pluto/datatype/request"

	"github.com/alecthomas/template"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	resp "github.com/leeif/pluto/datatype/response"
)

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
