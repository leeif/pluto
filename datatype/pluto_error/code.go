package pluto_error

import "net/http"

var (
	ServerError = 1000
	BadRequest  = NewPlutoError(http.StatusBadRequest, 1001, "Bad Request", nil)

	MailIsAlreadyRegister = NewPlutoError(http.StatusForbidden, 2001, "Mail is already been registered", nil)
	MailIsNotExsit        = NewPlutoError(http.StatusForbidden, 2002, "Mail is not exist", nil)
	MailIsNotVerified     = NewPlutoError(http.StatusForbidden, 2003, "Mail is not verified", nil)

	InvalidPassword        = NewPlutoError(http.StatusForbidden, 3001, "Invalid Password", nil)
	InvalidRefreshToken    = NewPlutoError(http.StatusForbidden, 3002, "Invalid Refresh Token", nil)
	InvalidMailVerifyToekn = NewPlutoError(http.StatusForbidden, 3003, "Invalid Mail Verify Token", nil)
	MailAlreadyVerified    = NewPlutoError(http.StatusBadRequest, 3004, "Mail is already verified", nil)
)

func NewServerError(err error) *PlutoError {
	return NewPlutoError(http.StatusInternalServerError, ServerError, "Server Error", err)
}
