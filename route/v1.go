package route

import (
	"path"

	"github.com/leeif/pluto/middleware"
)

func (r *Router) registerLoginV1Routes(prefix string) {
	routes := []PlutoRoute{
		{
			path:        "/register",
			description: "Register user with email",
			method:      "POST",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.Register,
		},
		{
			path:        "/register/verify/mail",
			description: "Send registration verification mail",
			method:      "POST",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.VerifyMail,
		},
		{
			path:        "/login",
			description: "Login with email or username",
			method:      "POST",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.Login,
		},
		{
			path:        "/login/google/mobile",
			description: "Login with google account for mobile app",
			method:      "POST",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.GoogleLoginMobile,
		},
		{
			path:        "/login/apple/mobile",
			description: "Login with apple account for mobile app",
			method:      "POST",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.AppleLoginMobile,
		},
		{
			path:        "/login/wechat/mobile",
			description: "Login with wechat account for mobile app",
			method:      "POST",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.WechatLoginMobile,
		},
		{
			path:        "/password/reset/mail",
			description: "Send password reset mail",
			method:      "POST",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.PasswordResetMail,
		},
	}
	r.registerRoutes(routes, path.Join(prefix, "/"), false)
}

func (r *Router) registerWebV1Routes(prefix string) {
	routes := []PlutoRoute{
		{
			path:        "/mail/verify/{token}",
			description: "Verify the mail registration",
			method:      "GET",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.RegistrationVerifyPage,
		},
		{
			path:        "/password/reset/{token}",
			description: "Reset password page",
			method:      "GET",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.ResetPasswordPage,
		},
		{
			path:        "/password/reset/{token}",
			description: "Reset password",
			method:      "POST",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.ResetPassword,
		},
		{
			path:        "/authorize",
			description: "Oauth authorize page",
			method:      "GET",
			middle:      middleware.AccessTokenAuthMiddleware,
			handler:     r.v1.AuthorizePage,
		},
		{
			path:        "/authorize",
			description: "Oauth authorize",
			method:      "POST",
			middle:      middleware.AccessTokenAuthMiddleware,
			handler:     r.v1.OAuthAuthorize,
		},
		{
			path:        "/login",
			description: "Login page",
			method:      "GET",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.LoginPage,
		},
		{
			path:        "/login",
			description: "Login with email or username",
			method:      "POST",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.OAuthLogin,
		},
	}
	r.registerRoutes(routes, path.Join(prefix, "/web"), true)
}

func (r *Router) registerHealthV1Routes(prefix string) {
	routes := []PlutoRoute{
		{
			path:        "/healthcheck",
			description: "Health check",
			method:      "GET",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.HealthCheck,
		},
	}
	r.registerRoutes(routes, path.Join(prefix, "/"), false)
}

func (r *Router) registerRBACV1Routes(prefix string) {
	routes := []PlutoRoute{
		{
			path:        "/rbac/role/create",
			description: "Create role",
			method:      "POST",
			middle:      middleware.AdminAuthMiddleware,
			handler:     r.v1.CreateRole,
		},
		{
			path:        "/rbac/scope/create",
			description: "Create scope",
			method:      "POST",
			middle:      middleware.AdminAuthMiddleware,
			handler:     r.v1.CreateScope,
		},
		{
			path:        "/rbac/role/scope",
			description: "Update scopes of the role",
			method:      "PUT",
			middle:      middleware.AdminAuthMiddleware,
			handler:     r.v1.RoleScopeUpdate,
		},
		{
			path:        "/rbac/role/scope/default",
			description: "Set the default scope of the role",
			method:      "PUT",
			middle:      middleware.AdminAuthMiddleware,
			handler:     r.v1.RoleDefaultScope,
		},
		{
			path:        "/rbac/application/create",
			description: "Create application",
			method:      "POST",
			middle:      middleware.AdminAuthMiddleware,
			handler:     r.v1.CreateApplication,
		},
		{
			path:        "/rbac/application/role/default",
			description: "Set the default role of the application",
			method:      "POST",
			middle:      middleware.AdminAuthMiddleware,
			handler:     r.v1.ApplicationDefaultRole,
		},
		{
			path:        "/rbac/application/list",
			description: "List all the applications",
			method:      "GET",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.ListApplications,
		},
		{
			path:        "/rbac/role/list",
			description: "List all the roles in the application",
			method:      "GET",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.ListRoles,
		},
		{
			path:        "/rbac/scope/list",
			description: "List all the scopes in the application",
			method:      "GET",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.ListScopes,
		},
		{
			path:        "/rbac/user/application/role",
			description: "Set the role of a user in application",
			method:      "POST",
			middle:      middleware.AdminAuthMiddleware,
			handler:     r.v1.SetUserRole,
		},
	}
	r.registerRoutes(routes, path.Join(prefix, "/rbac"), false)
}

func (r *Router) registerUserV1Routers(prefix string) {
	routes := []PlutoRoute{
		{
			path:        "/search",
			description: "Search the user using name or mail",
			method:      "GET",
			middle:      middleware.AdminAuthMiddleware,
			handler:     r.v1.FindUser,
		},
		{
			path:        "/count",
			description: "Get the count of the total users",
			method:      "GET",
			middle:      middleware.AdminAuthMiddleware,
			handler:     r.v1.UsersCount,
		},
		{
			path:        "/info/{userID}",
			description: "Get user info",
			method:      "GET",
			middle:      middleware.AccessTokenAuthMiddleware,
			handler:     r.v1.UserInfo,
		},
		{
			path:        "/info/update",
			description: "Update user info",
			method:      "POST",
			middle:      middleware.AccessTokenAuthMiddleware,
			handler:     r.v1.UpdateUserInfo,
		},
	}
	r.registerRoutes(routes, path.Join(prefix, "/user"), false)
}

func (r *Router) registerTokenV1Routes(prefix string) {
	routes := []PlutoRoute{
		{
			path:        "/refresh",
			description: "Refresh access token",
			method:      "POST",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.RefreshToken,
		},
		{
			path:        "/publickey",
			description: "Get the rsa public key",
			method:      "GET",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.PublicKey,
		},
		{
			path:        "/access/verify",
			description: "Verify access token",
			method:      "GET",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.VerifyAccessToken,
		},
	}
	r.registerRoutes(routes, path.Join(prefix, "/token"), false)
}

func (r *Router) registerOauthV1Routes(prefix string) {
	routes := []PlutoRoute{
		{
			path:        "/tokens",
			description: "request access token",
			method:      "POST",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.OAuthTokens,
		},
		{
			path:        "/client",
			description: "new client",
			method:      "POST",
			middle:      middleware.AccessTokenAuthMiddleware,
			handler:     r.v1.OAuthCreateClient,
		},
		{
			path:        "/client/approve",
			description: "approve client",
			method:      "PUT",
			middle:      middleware.AdminAuthMiddleware,
			handler:     r.v1.OAuthApproveClient,
		},
		{
			path:        "/client/deny",
			description: "deny client",
			method:      "PUT",
			middle:      middleware.AdminAuthMiddleware,
			handler:     r.v1.OAuthDenyClient,
		},
	}
	r.registerRoutes(routes, path.Join(prefix, "/oauth"), false)
}
