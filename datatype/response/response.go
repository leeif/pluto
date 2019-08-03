package datatype

const (
	STATUSOK    = "ok"
	STATUSERROR = "error"
)

type ReponseOK struct {
	Status string      `json:"status"`
	Body   interface{} `json:"body"`
}

type ReponseError struct {
	Status string      `json:"status"`
	Error  interface{} `json:"error"`
}
