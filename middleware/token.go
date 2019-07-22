package middleware

import "net/http"

func JWTAuthorization(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
}
