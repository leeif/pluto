package route

import (
	"net/http"

	"github.com/leeif/pluto/config"
	"github.com/leeif/pluto/middleware"

	"github.com/gorilla/mux"
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
	handler     negroni.HandlerFunc
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
	}
	r.registerRoutes(routes, "/api")
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
	r.registerRoutes(routes, "/")
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
	r.registerRoutes(routes, "/")
}

func (r *Router) registerRBACRoutes() {
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
			path:        "/rbac/scope/attach",
			description: "attach scope to role",
			method:      "POST",
			middle:      r.mw.AdminAuthMiddleware,
			handler:     r.AttachScope,
		},
		{
			path:        "/rbac/scope/detach",
			description: "detach scope to role",
			method:      "POST",
			middle:      r.mw.AdminAuthMiddleware,
			handler:     r.DetachScope,
		},
		{
			path:        "/rbac/role/scope/default",
			description: "set the default scope of the role",
			method:      "POST",
			middle:      r.mw.AdminAuthMiddleware,
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
	}
	r.registerRoutes(routes, "/api")
}

func (router *Router) registerRoutes(routes []PlutoRoute, prefix string) {
	sub := router.mux.PathPrefix(prefix).Subrouter()
	for _, r := range routes {
		sub.Handle(r.path, r.middle(r.handler)).Methods(r.method)
	}
}

func (r *Router) Register() {
	r.registerAPIRoutes()
	r.registerWebRoutes()
	r.registerHealthRoutes()
	r.registerRBACRoutes()
}

func NewRouter(mux *mux.Router, manager *manage.Manager, config *config.Config, logger *log.PlutoLog) *Router {
	return &Router{
		manager: manager,
		config:  config,
		logger:  logger,
		mw:      middleware.NewMiddle(logger),
		mux:     mux,
	}
}
