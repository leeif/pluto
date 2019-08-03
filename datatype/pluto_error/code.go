package pluto_error

import "net/http"

var (
	ServerError           = 1000
	BadRequest            = NewPlutoError(http.StatusBadRequest, 1001, "Bad Request", nil)
	MailIsAlreadyRegister = NewPlutoError(http.StatusForbidden, 1002, "Mail is already been registered", nil)
	MailIsNotExsit        = NewPlutoError(http.StatusForbidden, 1003, "Mail is not exist", nil)
	InvalidPassword       = NewPlutoError(http.StatusForbidden, 1004, "Invalid Password", nil)
	InvalidRefreshToken   = NewPlutoError(http.StatusForbidden, 1005, "Invalid Refresh Token", nil)
)

func NewServerError(err error) *PlutoError {
	return NewPlutoError(http.StatusInternalServerError, ServerError, "Server Error", err)
}
