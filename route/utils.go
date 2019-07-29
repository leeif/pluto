package route

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/leeif/pluto/datatype"
	resp "github.com/leeif/pluto/datatype/response"
)

func getBody(r *http.Request, revicer interface{}) *datatype.PlutoError {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return datatype.NewPlutoError(datatype.ServerError,
			errors.New("read body failed: "+err.Error()))
	}

	contentType := r.Header.Get("Content-type")
	if contentType == "application/json" {
		err := json.Unmarshal(body, &revicer)
		if err != nil {
			return datatype.NewPlutoError(datatype.ReqError,
				errors.New("Invalid JSON"))
		}
	}
	return nil
}

func responseOK(body interface{}, w http.ResponseWriter) error {
	response := resp.ReponseOK{}
	response.Status = resp.STATUSOK
	response.Body = body
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

func responseError(plutoError *datatype.PlutoError, w http.ResponseWriter) error {
	response := resp.ReponseError{}
	response.Status = resp.STATUSERROR
	switch plutoError.Type {
	case datatype.ReqError:
		response.Error = plutoError.Err.Error()
		w.WriteHeader(http.StatusBadRequest)
	case datatype.ServerError:
		response.Error = "server error"
		w.WriteHeader(http.StatusInternalServerError)
	}
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
