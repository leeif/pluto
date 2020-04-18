package config

type TokenConfig struct {
	AccessTokenExpire         int64 `kiper_value:"name:access_token_expire;help:expire time(s) of access token;default:3600"`
	ResetPasswordTokenExpire  int64 `kiper_value:"name:reset_password_token_expire;help:expire time(s) of reset password token;default:36000"`
	RegisterVerifyTokenExpire int64 `kiper_value:"name:register_verify_token_expire;help:expire time(s) of reset password result token;default:36000"`
	RefreshTokenExpire        int64 `kiper_value:"name:refresh_token_expire;help:expire time(s) of refresh token;default:604800"`
}

func newTokenConfig() *TokenConfig {
	return &TokenConfig{}
}
