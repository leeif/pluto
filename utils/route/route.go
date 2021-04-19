package route

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"

	"github.com/MuShare/pluto/config"

	"github.com/gorilla/context"
	"github.com/gorilla/schema"
	"github.com/wxnacy/wgo/arrays"

	perror "github.com/MuShare/pluto/datatype/pluto_error"
	"github.com/MuShare/pluto/datatype/request"
	resp "github.com/MuShare/pluto/datatype/response"
	"github.com/MuShare/pluto/log"
	"github.com/MuShare/pluto/utils/jwt"
	"github.com/MuShare/pluto/utils/view"
)

func GetBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	baseURL := fmt.Sprintf("%s://%s", scheme, r.Host)

	// forwarded host
	forwardedHost := r.Header.Get("X-Forwarded-Host")
	if forwardedHost != "" {
		baseURL = forwardedHost
	}
	return baseURL
}

func GetRequestData(r *http.Request, receiver interface{}) *perror.PlutoError {

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	contentType := r.Header.Get("Content-type")
	if strings.Contains(contentType, "application/json") {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return perror.ServerError.Wrapper(errors.New("Read body failed: " + err.Error()))
		}
		err = json.Unmarshal(body, &receiver)
		if err != nil {
			return perror.BadRequest.Wrapper(err)
		}
	} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
		err := r.ParseForm()
		if err != nil {
			return perror.BadRequest
		}
		err = decoder.Decode(receiver, r.PostForm)
		if err != nil {
			return perror.BadRequest.Wrapper(err)
		}
	}

	// parse url parameter
	if err := decoder.Decode(receiver, r.URL.Query()); err != nil {
		return perror.BadRequest.Wrapper(err)
	}

	pr, ok := receiver.(request.PlutoRequest)
	// check request body validation
	if ok && !pr.Validation() {
		return perror.BadRequest
	}
	return nil
}

func GetAccessPayload(r *http.Request) (*jwt.AccessPayload, *perror.PlutoError) {
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

func ResponseOK(body interface{}, w http.ResponseWriter) *perror.PlutoError {
	response := resp.ReponseOK{}
	response.Status = resp.STATUSOK
	response.Body = body
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(response)
	if err != nil {
		return perror.HTTPResponseError.Wrapper(err)
	}
	_, err = w.Write(b)
	if err != nil {
		return perror.HTTPResponseError.Wrapper(err)
	}
	return nil
}

func ResponseError(plutoError *perror.PlutoError, w http.ResponseWriter) *perror.PlutoError {
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
		return perror.HTTPResponseError.Wrapper(err)
	}
	_, err = w.Write(b)
	if err != nil {
		return perror.HTTPResponseError.Wrapper(err)
	}
	return nil
}

func ResponseHTMLOK(file string, data interface{}, r *http.Request, w http.ResponseWriter, config *config.Config) *perror.PlutoError {

	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(http.StatusOK)
	vw, err := view.GetView()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}
	t, err := vw.Parse(r.Header.Get("Accept-Language"), file)
	if err != nil {
		return perror.HTTPResponseError.Wrapper(err)
	}
	data, err = setFooter(data, config)
	if err != nil {
		return perror.HTTPResponseError.Wrapper(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		return perror.HTTPResponseError.Wrapper(err)
	}
	return nil
}

func ResponseHTMLError(file string, data interface{}, r *http.Request, w http.ResponseWriter, status int, config *config.Config) *perror.PlutoError {
	w.Header().Set("Content-type", "text/html")
	w.WriteHeader(status)
	vw, err := view.GetView()
	if err != nil {
		return perror.ServerError.Wrapper(err)
	}
	t, err := vw.Parse(r.Header.Get("Accept-Language"), file)
	if err != nil {
		return perror.HTTPResponseError.Wrapper(err)
	}
	if err != nil {
		return perror.HTTPResponseError.Wrapper(err)
	}
	data, err = setFooter(data, config)
	if err != nil {
		return perror.HTTPResponseError.Wrapper(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		return perror.HTTPResponseError.Wrapper(err)
	}
	return nil
}

func setFooter(data interface{}, config *config.Config) (d map[string]interface{}, err error) {
	var dataMap map[string]interface{}
	if data == nil {
		dataMap = make(map[string]interface{})
	} else {
		inrec, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		json.Unmarshal(inrec, &dataMap)
	}
	footer := "<div class=\"footer col-12 text-center fixed-bottom\"><div class=\"link\" style=\"margin: 1em 0\">" + config.Server.HTMLFooter + "</div></div>"
	dataMap["Footer"] = footer
	return dataMap, err
}

// Redirects to a new path while keeping current request's query string
func RedirectWithQueryString(to string, query url.Values, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s%s", to, GetQueryString(query)), http.StatusFound)
}

// Redirects to a new path with the query string moved to the URL fragment
func RedirectWithFragment(to string, query url.Values, w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, fmt.Sprintf("%s#%s", to, query.Encode()), http.StatusFound)
}

// Helper function to handle redirecting failed or declined authorization
func ErrorRedirect(w http.ResponseWriter, r *http.Request, redirectURI *url.URL, errCode int, state, responseType string) {
	query := redirectURI.Query()
	query.Set("error_code", strconv.Itoa(errCode))
	if state != "" {
		query.Set("state", state)
	}
	if responseType == "code" {
		RedirectWithQueryString(redirectURI.String(), query, w, r)
	}
	if responseType == "token" {
		RedirectWithFragment(redirectURI.String(), query, w, r)
	}
}

// Returns string encoded query string of the request
func GetQueryString(query url.Values, except ...string) string {
	if len(query) == 0 {
		return ""
	}
	res := make([]string, 0)
	for k, v := range query {
		if arrays.Contains(except, k) != -1 {
			continue
		}
		value := ""
		if len(v) > 0 {
			value = v[0]
		}
		res = append(res, fmt.Sprintf("%s=%s", k, value))
	}
	return "?" + strings.Join(res, "&")
}

func PlutoLog(logger *log.PlutoLog, pe *perror.PlutoError, r *http.Request) {
	url := r.URL.String()
	if pe.LogError != nil {
		logger.Error(fmt.Sprintf("[(%s)%s]:%s", r.Method, url, pe.LogError.Error()))
		fmt.Println(string(debug.Stack()))
	}
	if pe.HTTPError != nil {
		logger.Debug(fmt.Sprintf("[(%s)%s]:%s", r.Method, url, pe.HTTPError.Error()))
	}
}
