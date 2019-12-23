package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/leeif/pluto/utils/jwt"

	"github.com/gorilla/context"
	"github.com/gorilla/schema"
	"github.com/leeif/pluto/datatype/request"
	"github.com/leeif/pluto/utils/view"

	perror "github.com/leeif/pluto/datatype/pluto_error"
	resp "github.com/leeif/pluto/datatype/response"
)

func getBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func getBody(r *http.Request, reciver interface{}) *perror.PlutoError {

	contentType := r.Header.Get("Content-type")
	if strings.Contains(contentType, "application/json") {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return perror.ServerError.Wrapper(errors.New("Read body failed: " + err.Error()))
		}
		err = json.Unmarshal(body, &reciver)
		if err != nil {
			return perror.BadRequest
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		err := r.ParseForm()
		if err != nil {
			return perror.BadRequest
		}
		decoder := schema.NewDecoder()
		err = decoder.Decode(reciver, r.PostForm)
		if err != nil {
			return perror.BadRequest.Wrapper(err)
		}
	}
	pr, ok := reciver.(request.PlutoRequest)
	// check request body validation
	if ok && !pr.Validation() {
		return perror.BadRequest
	}
	return nil
}

func getQuery(r *http.Request, reciver interface{}) *perror.PlutoError {
	decoder := schema.NewDecoder()
	if err := decoder.Decode(reciver, r.URL.Query()); err != nil {
		return perror.BadRequest.Wrapper(err)
	}
	pr, ok := reciver.(request.PlutoRequest)
	// check request body validation
	if ok && !pr.Validation() {
		return perror.BadRequest
	}

	return nil
}

func getAccessPayload(r *http.Request) (*jwt.AccessPayload, *perror.PlutoError) {
	perr := context.Get(r, "pluto_error")
	if perr != nil {
		err, ok := perr.(*perror.PlutoError)
		if !ok {
			err = perror.ServerError.Wrapper(fmt.Errorf("Unknow error"))
		}
		return nil, err
	}

	accessPayload := context.Get(r, "payload")

	if accessPayload == nil {
		err := perror.ServerError.Wrapper(fmt.Errorf("Access token payload is empty"))
		return nil, err
	}

	payload, ok := accessPayload.(*jwt.AccessPayload)
	if !ok {
		err := perror.ServerError.Wrapper(fmt.Errorf("Not a access token payload"))
		return nil, err
	}
	return payload, nil
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
	t, err := view.Parse(file)
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
	t, err := view.Parse(file)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	err = t.Execute(w, data)
	if err != nil {
		return err
	}
	return nil
}
