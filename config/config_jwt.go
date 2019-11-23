package config

type JWTConfig struct {
	AccessTokenExpire         int64 `kiper_value:"name:access_token_expire;help:expire time(s) of access token;default:3600"`
	ResetPasswordTokenExpire  int64 `kiper_value:"name:reset_password_token_expire;help:expire time(s) of reset password token;default:36000"`
	RegisterVerifyTokenExpire int64 `kiper_value:"name:register_verify_token_expire;help:expire time(s) of reset password result token;default:36000"`
}

func newJWTConfig() *JWTConfig {
	return &JWTConfig{}
}
