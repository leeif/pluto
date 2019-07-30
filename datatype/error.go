package datatype

const (
	ServerError = iota
	ReqError
)

type PlutoError struct {
	Type int
	Err  error
}

func NewPlutoError(t int, err error) *PlutoError {
	return &PlutoError{
		Type: t,
		Err:  err,
	}
}
