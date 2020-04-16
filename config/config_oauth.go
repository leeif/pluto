package config

type OAuthConfig struct {
	AuthorizeCodeExpire int `kiper_value:"name:authorize_code_expire;help:expire time(s) of authorize;default:3600"`
}

func newOAuthConfig() *OAuthConfig {
	return &OAuthConfig{}
}
