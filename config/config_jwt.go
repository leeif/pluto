package config

type JWTConfig struct {
	AccessTokenExpire              int64 `kiper_value:"name:access_token_expire;help:expire time(s) of access token;default:600"`
	ResetPasswordTokenExpire       int64 `kiper_value:"name:reset_password_token;help:expire time(s) of reset password token;default:1200"`
	ResetPasswordResultTokenExpire int64 `kiper_value:"name:reset_password_result_token;help:expire time(s) of reset password result token;default:300"`
	RegisterVerifyTokenExpire      int64 `kiper_value:"name:register_verify_token;help:expire time(s) of reset password result token;default:300"`
}

func newJWTConfig() *JWTConfig {
	return &JWTConfig{}
}
