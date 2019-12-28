package pluto_error

import (
	"errors"
	"fmt"
)

type PlutoError struct {
	HTTPCode  int
	HTTPError error
	PlutoCode int
	LogError  error
}

func (pe *PlutoError) Wrapper(err error) *PlutoError {
	pe.LogError = fmt.Errorf("%s\n%w", err.Error(), pe.LogError)
	return pe
}

func NewPlutoError(httpCode int, plutoCode int, httpError string, logError error) *PlutoError {
	return &PlutoError{
		HTTPCode:  httpCode,
		HTTPError: errors.New(httpError),
		PlutoCode: plutoCode,
		LogError:  logError,
	}
}
