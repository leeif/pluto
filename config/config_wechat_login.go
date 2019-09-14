package config

type WechatLoginConfig struct {
	AppID  *string `kiper_value:"name:app_id;help:wechat app id"`
	Secret *string `kiper_value:"name:secret;help:wechat secret"`
}

func newWechatLoginConfig() *WechatLoginConfig {
	return &WechatLoginConfig{}
}
