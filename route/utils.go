package route

import (
	"errors"
	"net/http"
	"net/url"
	"path"

	"github.com/urfave/negroni"

	routeUtils "github.com/leeif/pluto/utils/route"

	perror "github.com/leeif/pluto/datatype/pluto_error"
)

func (router *Router) plutoHandlerWrapper(handler func(http.ResponseWriter, *http.Request) *perror.PlutoError) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		defer func() {
			if err := recover(); err != nil {
				var perr *perror.PlutoError
				switch err.(type) {
				case error:
					perr = perror.ServerError.Wrapper(err.(error))
				default:
					perr = perror.ServerError.Wrapper(errors.New("unknown error"))
				}
				routeUtils.ResponseError(perr, w)
				routeUtils.PlutoLog(router.logger, perr, r)
			}
		}()

		if err := handler(w, r); err != nil {
			if err.PlutoCode != perror.HTTPResponseError.PlutoCode {
				routeUtils.ResponseError(err, w)
			}
			routeUtils.PlutoLog(router.logger, err, r)
			return
		}

		next(w, r)
	}
}

func (router *Router) plutoWebHandlerWrapper(handler func(http.ResponseWriter, *http.Request) *perror.PlutoError) negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {

		type Data struct {
			Error *perror.PlutoError
		}

		defer func() {
			if err := recover(); err != nil {
				var perr *perror.PlutoError
				switch err.(type) {
				case error:
					perr = perror.ServerError.Wrapper(err.(error))
				default:
					perr = perror.ServerError.Wrapper(errors.New("unknown error"))
				}
				data := &Data{
					Error: perr,
				}
				routeUtils.ResponseHTMLError("error.html", data, r, w, http.StatusInternalServerError)
				routeUtils.PlutoLog(router.logger, perr, r)
			}
		}()

		if err := handler(w, r); err != nil {
			if err.PlutoCode == perror.Unauthorized.PlutoCode {
				u, err := url.Parse(routeUtils.GetBaseURL(r))
				if err != nil {
					return
				}
				u.Path = path.Join(u.Path, "/web/login")
				query := r.URL.Query()
				query.Set("login_redirect_uri", r.URL.Path)
				routeUtils.RedirectWithQueryString(u.String(), query, w, r)
			} else if err.PlutoCode != perror.HTTPResponseError.PlutoCode {
				data := &Data{
					Error: err,
				}
				routeUtils.ResponseHTMLError("error.html", data, r, w, http.StatusInternalServerError)
			}
			routeUtils.PlutoLog(router.logger, err, r)
			return
		}

		next(w, r)
	}
}
