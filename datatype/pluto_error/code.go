package pluto_error

import "net/http"

var (
	ServerError       = NewPlutoError(http.StatusInternalServerError, 1000, "Server Error", nil)
	BadRequest        = NewPlutoError(http.StatusBadRequest, 1001, "Bad Request", nil)
	Forbidden         = NewPlutoError(http.StatusForbidden, 1002, "Permission Forbidden", nil)
	Unauthorized      = NewPlutoError(http.StatusUnauthorized, 1003, "Request is unauthorized", nil)
	HTTPResponseError = NewPlutoError(http.StatusInternalServerError, 1004, "http writer error", nil)
	NotFound          = NewPlutoError(http.StatusNotFound, 1005, "Path not found", nil)

	MailIsAlreadyRegister = NewPlutoError(http.StatusForbidden, 2001, "Mail is already been registered", nil)
	MailNotExist          = NewPlutoError(http.StatusForbidden, 2002, "Mail does not exist", nil)
	MailIsNotVerified     = NewPlutoError(http.StatusForbidden, 2003, "Mail is not verified", nil)
	MailAlreadyVerified   = NewPlutoError(http.StatusBadRequest, 2004, "Mail is already verified", nil)
	UsernameNotExist      = NewPlutoError(http.StatusForbidden, 2005, "Username does not exist", nil)

	InvalidPassword      = NewPlutoError(http.StatusForbidden, 3001, "Invalid Password", nil)
	InvalidRefreshToken  = NewPlutoError(http.StatusForbidden, 3002, "Invalid Refresh Token", nil)
	InvalidJWTToken      = NewPlutoError(http.StatusForbidden, 3003, "Invalid JWT Token", nil)
	InvalidGoogleIDToken = NewPlutoError(http.StatusForbidden, 3004, "Invalid Google ID Token", nil)
	InvalidWechatCode    = NewPlutoError(http.StatusForbidden, 3005, "Invalid Wechat Code", nil)
	InvalidAvatarFormat  = NewPlutoError(http.StatusBadRequest, 3006, "Invalid Avatar Format", nil)
	InvalidAppleIDToken  = NewPlutoError(http.StatusForbidden, 3007, "Invalid Apple ID Token", nil)
	JWTTokenExpired      = NewPlutoError(http.StatusForbidden, 3008, "JWT Token Expired", nil)
	InvalidAccessToken   = NewPlutoError(http.StatusForbidden, 3009, "Invalid Access Token", nil)
	InvalidApplication   = NewPlutoError(http.StatusForbidden, 3010, "Invalid Application", nil)

	ScopeExists         = NewPlutoError(http.StatusForbidden, 4001, "Scope already exists", nil)
	ScopeNotExist       = NewPlutoError(http.StatusForbidden, 4002, "Scope not exist", nil)
	ScopeAttached       = NewPlutoError(http.StatusForbidden, 4003, "Scope already attached", nil)
	ApplicationExists   = NewPlutoError(http.StatusForbidden, 4004, "Application already exists", nil)
	ApplicationNotExist = NewPlutoError(http.StatusForbidden, 4005, "Application does not exist", nil)
	RoleExists          = NewPlutoError(http.StatusForbidden, 4006, "Role does not exist", nil)
	RoleNotExist        = NewPlutoError(http.StatusForbidden, 4007, "Role already exists", nil)
	NotPlutoAdmin       = NewPlutoError(http.StatusForbidden, 4008, "Not the pluto admin", nil)
	UserNotExist        = NewPlutoError(http.StatusForbidden, 4009, "User is not exist", nil)

	OAuthInvalidGrantType          = NewPlutoError(http.StatusForbidden, 5001, "Invalid grant type", nil)
	OAuthAuthorizationCodeNotFound = NewPlutoError(http.StatusForbidden, 5002, "Authorization code not found", nil)
	OAuthInvalidRedirectURI        = NewPlutoError(http.StatusForbidden, 5003, "Invlid redirect uri", nil)
	OAuthAuthorizationCodeExpired  = NewPlutoError(http.StatusForbidden, 5004, "Authorization code expired", nil)
	OAuthClientIDOrSecretNotFound  = NewPlutoError(http.StatusUnauthorized, 5005, "client_id or secret not found", nil)
	OAuthClientExist               = NewPlutoError(http.StatusUnauthorized, 5006, "Client already exists", nil)
	OAuthInvalidClient             = NewPlutoError(http.StatusForbidden, 5007, "Invalid client", nil)
)
