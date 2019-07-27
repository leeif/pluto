package datatype

const (
	STATUSOK    = "ok"
	STATUSERROR = "error"
)

type Reponse struct {
	Status string      `json:"status"`
	Body   interface{} `json:"body,omitempty"`
}
