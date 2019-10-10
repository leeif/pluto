package pluto_error

import "errors"

type PlutoError struct {
	HTTPCode  int
	HTTPError error
	PlutoCode int
	LogError  error
}

func (pe *PlutoError) Wrapper(err error) *PlutoError {
	str := ""
	if pe.LogError != nil {
		str += pe.LogError.Error() + "\n"
	}
	if err != nil {
		str += err.Error() + "\n"
	}
	newPE := PlutoError(*pe)
	newPE.LogError = errors.New(str)
	return &newPE
}

func NewPlutoError(httpCode int, plutoCode int, httpError string, logError error) *PlutoError {
	return &PlutoError{
		HTTPCode:  httpCode,
		HTTPError: errors.New(httpError),
		PlutoCode: plutoCode,
		LogError:  logError,
	}
}
