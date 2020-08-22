package route

import (
	"path"

	"github.com/leeif/pluto/middleware"
)

func (r *Router) registerUserV1Routes(prefix string) {
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
			path:        "/login/account",
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
			path:        "/binding",
			description: "Bind mail, google, wechat, apple account",
			method:      "POST",
			middle:      middleware.AccessTokenAuthMiddleware,
			handler:     r.v1.Binding,
		},
		{
			path:        "/unbinding",
			description: "Unbind mail, google, wechat, apple account",
			method:      "POST",
			middle:      middleware.AccessTokenAuthMiddleware,
			handler:     r.v1.Unbinding,
		},
		{
			path:        "/password/reset/mail",
			description: "Send password reset mail",
			method:      "POST",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.PasswordResetMail,
		},
		{
			path:        "/search",
			description: "Search the user using name or mail",
			method:      "GET",
			middle:      middleware.PlutoAdminAuthMiddleware,
			handler:     r.v1.FindUser,
		},
		{
			path:        "/summary",
			description: "Get the summary of users",
			method:      "GET",
			middle:      middleware.PlutoAdminAuthMiddleware,
			handler:     r.v1.UserSummary,
		},
		{
			path:        "/info",
			description: "Get user info",
			method:      "GET",
			middle:      middleware.AccessTokenAuthMiddleware,
			handler:     r.v1.UserInfo,
		},
		{
			path:        "/info",
			description: "Update user info",
			method:      "PUT",
			middle:      middleware.AccessTokenAuthMiddleware,
			handler:     r.v1.UpdateUserInfo,
		},
		{
			path:        "/info/public",
			description: "Get user public info like avatar, name",
			method:      "GET",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.PublicUserInfo,
		},
	}
	r.registerRoutes(routes, path.Join(prefix, "/user"), false)
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
			path:        "/role/create",
			description: "Create role",
			method:      "POST",
			middle:      middleware.PlutoAdminAuthMiddleware,
			handler:     r.v1.CreateRole,
		},
		{
			path:        "/scope/create",
			description: "Create scope",
			method:      "POST",
			middle:      middleware.PlutoAdminAuthMiddleware,
			handler:     r.v1.CreateScope,
		},
		{
			path:        "/role/scope",
			description: "Update scopes of the role",
			method:      "PUT",
			middle:      middleware.PlutoAdminAuthMiddleware,
			handler:     r.v1.RoleScopeUpdate,
		},
		{
			path:        "/role/scope/default",
			description: "Set the default scope of the role",
			method:      "PUT",
			middle:      middleware.PlutoAdminAuthMiddleware,
			handler:     r.v1.RoleDefaultScope,
		},
		{
			path:        "/application/create",
			description: "Create application",
			method:      "POST",
			middle:      middleware.PlutoAdminAuthMiddleware,
			handler:     r.v1.CreateApplication,
		},
		{
			path:        "/application/role/default",
			description: "Set the default role of the application",
			method:      "POST",
			middle:      middleware.PlutoAdminAuthMiddleware,
			handler:     r.v1.ApplicationDefaultRole,
		},
		{
			path:        "/application/list",
			description: "List all the applications",
			method:      "GET",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.ListApplications,
		},
		{
			path:        "/role/list",
			description: "List all the roles in the application",
			method:      "GET",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.ListRoles,
		},
		{
			path:        "/scope/list",
			description: "List all the scopes in the application",
			method:      "GET",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.ListScopes,
		},
		{
			path:        "/user/application/role",
			description: "Set the role of a user in application",
			method:      "POST",
			middle:      middleware.PlutoAdminAuthMiddleware,
			handler:     r.v1.SetUserRole,
		},
	}
	r.registerRoutes(routes, path.Join(prefix, "/rbac"), false)
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
			description: "Request access token",
			method:      "POST",
			middle:      middleware.NoAuthMiddleware,
			handler:     r.v1.OAuthTokens,
		},
		{
			path:        "/client",
			description: "Get clients",
			method:      "GET",
			middle:      middleware.PlutoUserAuthMiddleware,
			handler:     r.v1.OAuthGetClient,
		},
		{
			path:        "/client",
			description: "Create client",
			method:      "POST",
			middle:      middleware.PlutoUserAuthMiddleware,
			handler:     r.v1.OAuthCreateClient,
		},
		{
			path:        "/client/status",
			description: "Change the client status",
			method:      "PUT",
			middle:      middleware.PlutoAdminAuthMiddleware,
			handler:     r.v1.OAuthApproveClient,
		},
	}
	r.registerRoutes(routes, path.Join(prefix, "/oauth"), false)
}
