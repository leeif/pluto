package route

import (
	"net/http"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/middleware"

	"github.com/gorilla/mux"
	perror "github.com/leeif/pluto/datatype/pluto_error"
	"github.com/leeif/pluto/log"
	"github.com/leeif/pluto/manage"

	"github.com/urfave/negroni"
)

type middle func(handlers ...negroni.HandlerFunc) http.Handler

type PlutoRoute struct {
	path        string
	description string
	method      string
	middle      middle
	handler     func(w http.ResponseWriter, r *http.Request) *perror.PlutoError
}

type Router struct {
	manager *manage.Manager
	config  *config.Config
	logger  *log.PlutoLog
	mw      *middleware.Middleware
	mux     *mux.Router
}

func (r *Router) registerAPIRoutes() {
	routes := []PlutoRoute{
		{
			path:        "/user/register",
			description: "register user with email",
			method:      "POST",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.register,
		},
		{
			path:        "/user/register/verify/mail",
			description: "send registration verification mail",
			method:      "POST",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.verifyMail,
		},
		{
			path:        "/user/login",
			description: "login with mail",
			method:      "POST",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.login,
		},
		{
			path:        "/user/login/google/mobile",
			description: "login with google account for mobile app",
			method:      "POST",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.googleLoginMobile,
		},
		{
			path:        "/user/login/apple/mobile",
			description: "login with apple account for mobile app",
			method:      "POST",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.appleLoginMobile,
		},
		{
			path:        "/user/login/wechat/mobile",
			description: "login with wechat account for mobile app",
			method:      "POST",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.wechatLoginMobile,
		},
		{
			path:        "/user/password/reset/mail",
			description: "send password reset mail",
			method:      "POST",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.passwordResetMail,
		},
		{
			path:        "/user/info/me",
			description: "get user info",
			method:      "GET",
			middle:      r.mw.AccessTokenAuthMiddleware,
			handler:     r.userInfo,
		},
		{
			path:        "/user/info/me/update",
			description: "update user info",
			method:      "POST",
			middle:      r.mw.AccessTokenAuthMiddleware,
			handler:     r.updateUserInfo,
		},
		{
			path:        "/user/count",
			description: "get the count of the total users",
			method:      "GET",
			middle:      r.mw.AdminAuthMiddleware,
			handler:     r.UsersCount,
		},
		{
			path:        "/auth/refresh",
			description: "refresh access token",
			method:      "POST",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.refreshToken,
		},
		{
			path:        "/auth/publickey",
			description: "get the rsa public key",
			method:      "GET",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.publicKey,
		},
		{
			path:        "/auth/token/verify/access",
			description: "verify access token",
			method:      "POST",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.verifyAccessToken,
		},
	}
	r.registerRoutes(routes, "/api", false)
}

func (r *Router) registerWebRoutes() {
	routes := []PlutoRoute{
		{
			path:        "/mail/verify/{token}",
			description: "verify the mail registration",
			method:      "GET",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.registrationVerifyPage,
		},
		{
			path:        "/password/reset/{token}",
			description: "reset password page",
			method:      "GET",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.resetPasswordPage,
		},
		{
			path:        "/password/reset/{token}",
			description: "reset password",
			method:      "POST",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.resetPassword,
		},
	}
	r.registerRoutes(routes, "/", true)
}

func (r *Router) registerHealthRoutes() {
	routes := []PlutoRoute{
		{
			path:        "/healthcheck",
			description: "health check api",
			method:      "GET",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.healthCheck,
		},
	}
	r.registerRoutes(routes, "/", false)
}

func (r *Router) registerAdminRoutes() {
	routes := []PlutoRoute{
		{
			path:        "/rbac/role/create",
			description: "create role",
			method:      "POST",
			middle:      r.mw.AdminAuthMiddleware,
			handler:     r.CreateRole,
		},
		{
			path:        "/rbac/scope/create",
			description: "create scope",
			method:      "POST",
			middle:      r.mw.AdminAuthMiddleware,
			handler:     r.CreateScope,
		},
		{
			path:        "/rbac/role/scope/attach",
			description: "attach scope to role",
			method:      "POST",
			middle:      r.mw.AdminAuthMiddleware,
			handler:     r.AttachScope,
		},
		{
			path:        "/rbac/role/scope/detach",
			description: "detach scope to role",
			method:      "POST",
			middle:      r.mw.AdminAuthMiddleware,
			handler:     r.DetachScope,
		},
		{
			path:        "/rbac/role/scope/batch",
			description: "detach scope to role",
			method:      "POST",
			middle:      r.mw.AdminAuthMiddleware,
			handler:     r.RoleScopeBatchUpdate,
		},
		{
			path:        "/rbac/role/scope/default",
			description: "set the default scope of the role",
			method:      "POST",
			middle:      r.mw.AdminAuthMiddleware,
			handler:     r.RoleDefaultScope,
		},
		{
			path:        "/rbac/application/create",
			description: "create application",
			method:      "POST",
			middle:      r.mw.AdminAuthMiddleware,
			handler:     r.CreateApplication,
		},
		{
			path:        "/rbac/application/role/default",
			description: "set the default role of the application",
			method:      "POST",
			middle:      r.mw.AdminAuthMiddleware,
			handler:     r.ApplicationDefaultRole,
		},
		{
			path:        "/rbac/application/list",
			description: "list all the applications",
			method:      "GET",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.ListApplications,
		},
		{
			path:        "/rbac/role/list",
			description: "list all the roles in the application",
			method:      "GET",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.ListRoles,
		},
		{
			path:        "/rbac/scope/list",
			description: "list all the scopes in the application",
			method:      "GET",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.ListScopes,
		},
		{
			path:        "/find/user",
			description: "find the user through name or mail",
			method:      "GET",
			middle:      r.mw.AdminAuthMiddleware,
			handler:     r.FindUser,
		},
		{
			path:        "/rbac/user/application/role",
			description: "set the role of a user in application",
			method:      "POST",
			middle:      r.mw.AdminAuthMiddleware,
			handler:     r.SetUserRole,
		},
	}
	r.registerRoutes(routes, "/api", false)
}

func (r *Router) registerOauth2Routes() {
	routes := []PlutoRoute{
		{
			path:        "/tokens",
			description: "request access token",
			method:      "POST",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.OAuth2Tokens,
		},
		{
			path:        "/authorize",
			description: "authorize page",
			method:      "GET",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.OAuth2Authorize,
		},
		{
			path:        "/login",
			description: "login for oauth2",
			method:      "GET",
			middle:      r.mw.NoAuthMiddleware,
			handler:     r.OAuth2Login,
		},
	}
	r.registerRoutes(routes, "/oauth2", false)
}

func (router *Router) registerRoutes(routes []PlutoRoute, prefix string, isWeb bool) {
	sub := router.mux.PathPrefix(prefix).Subrouter()
	for _, r := range routes {
		// options method for cors
		if isWeb {
			sub.Handle(r.path, r.middle(router.plutoWebHandlerWrapper(r.handler))).Methods(r.method)
		} else {
			sub.Handle(r.path, r.middle(router.plutoHandlerWrapper(r.handler))).Methods(r.method)
		}
	}
}

func (router *Router) Register() {
	router.registerAPIRoutes()
	router.registerWebRoutes()
	router.registerHealthRoutes()
	router.registerAdminRoutes()
}

func NewRouter(mux *mux.Router, manager *manage.Manager, config *config.Config, logger *log.PlutoLog) *Router {
	return &Router{
		manager: manager,
		config:  config,
		logger:  logger,
		mw:      middleware.NewMiddle(logger, config),
		mux:     mux,
	}
}
