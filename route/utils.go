package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/urfave/negroni"
	"io/ioutil"
	"net/http"
	"net/url"
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

func (router *Router) plutoHandlerWrapper(handler func(http.ResponseWriter, *http.Request) *perror.PlutoError) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if err := handler(w, r); err != nil {
			responseError(err, w)
			router.plutoLog(err, r)
			return
		}

		next(w, r)
	}
}

func (router *Router) plutoWebHandlerWrapper(handler func(http.ResponseWriter, *http.Request) *perror.PlutoError) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if err := handler(w, r); err != nil {
			responseHTMLError("error.html", nil, w, http.StatusInternalServerError)
			router.plutoLog(err, r)
			return
		}

		next(w, r)
	}
}

func (router *Router) plutoLog(pe *perror.PlutoError, r *http.Request) {
	url := r.URL.String()
	if pe.LogError != nil {
		router.logger.Error(fmt.Sprintf("[%s]:%s", url, pe.LogError.Error()))
	}
	if pe.HTTPError != nil {
		router.logger.Debug(fmt.Sprintf("[%s]:%s", url, pe.HTTPError.Error()))
	}
}

// Redirects to a new path while keeping current request's query string
func redirectWithQueryString(to string, query url.Values, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s%s", to, getQueryString(query)), http.StatusFound)
}

// Redirects to a new path with the query string moved to the URL fragment
func redirectWithFragment(to string, query url.Values, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s#%s", to, query.Encode()), http.StatusFound)
}

// Returns string encoded query string of the request
func getQueryString(query url.Values) string {
	encoded := query.Encode()
	if len(encoded) > 0 {
		encoded = fmt.Sprintf("?%s", encoded)
	}
	return encoded
}
