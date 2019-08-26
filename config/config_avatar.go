package config

type AvatarConfig struct {
	Bucket          *string `kiper_value:"name:bucket;help:s3 bucket"`
	EndPoint        *string `kiper_value:"name:endpoint;help:end point"`
	AccessKeyID     *string `kiper_value:"name:accesskeyid;help:access key id"`
	AccessKeySecret *string `kiper_value:"name:accesskeysecret;help:access key secret"`
	CDN             *string `kiper_value:"name:cdn;help:cdn url"`
}

func newAvatarConfig() *AvatarConfig {
	return &AvatarConfig{}
}
