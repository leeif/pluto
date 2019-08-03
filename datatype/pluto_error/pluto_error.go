package pluto_error

import "errors"

type PlutoError struct {
	HTTPCode  int
	HTTPError error
	PlutoCode int
	LogError  error
}

func NewPlutoError(httpCode int, plutoCode int, httpError string, logError error) *PlutoError {
	return &PlutoError{
		HTTPCode:  httpCode,
		HTTPError: errors.New(httpError),
		PlutoCode: plutoCode,
		LogError:  logError,
	}
}
