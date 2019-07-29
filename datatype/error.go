package datatype

import "fmt"

const (
	ServerError = iota
	ReqError
)

type PlutoError struct {
	Type int
	Err  error
}

func NewPlutoError(t int, err error) *PlutoError {
	fmt.Println(err)
	return &PlutoError{
		Type: t,
		Err:  err,
	}
}
