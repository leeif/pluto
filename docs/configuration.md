# Configuration

Pluto use both command line flags and config files for configuration.

If both command line flags and config files are set, the config file's will be used.

Pluto use [kiper](https://github.com/leeif/kiper), which is a wrapper of kingpin and viper to manage the configrations.

## Use Flag

### server

|  Command line flag  |  Description  | Type | Default |
| ---- | ---- | ---- | ---- |
|--server.port| port which Pluto server listens | string | 8010 |
|--server.skip_register_verify_mail| skip mail verification when register with mail | bool | false |

### log

|  Command line flag  |  Description  | Type | Default |
| ---- | ---- | ---- | ---- |
|--log.level|log level (debug, info, warn, error)| string|info|
|--log.format|log format (json, logfmt)| string|logfmt|
|--log.file|log file path, use stdout if this setting is empty | string|""|

### rsa

RSA private/public files are used to sign/verify JWTs.

|  Command line flag  |  Description  | Type | Default |
| ---- | ---- | ---- | ---- |
|--rsa.path|rsa file path|string|./|
|--rsa.name|rsa file name|string|id_rsa|

### database

database now only support MySQL.

|  Command line flag  |  Description  | Type | Default |
| ---- | ---- | ---- | ---- |
|--database.type|database type, now mysql only|string|mysql|
|--database.host|database host|string|127.0.0.1|
|--database.user|database user|string|root|
|--database.password|database password|string|""|
|--database.db|db name|string|pluto|

### mail

Pluto send verification mail or reset password mail through smtp server. If the smtp server are not set, 500 status code will be return when request mail related APIs.

|  Command line flag  |  Description  | Type | Default |
| ---- | ---- | ---- | ---- |
|--mail.smtp|smtp server of mail|string|""|
|--mail.user|user of smtp server|string|""|
|--mail.password|password of smtp server|string|""|

### avatar

Pluto random set an avatar for user when register using [https://www.gravatar.com/avatar/](https://www.gravatar.com/avatar/)

You can also save avatar in remote object storage, now only support aliyun OSS.

|  Command line flag  |  Description  | Type | Default |
| ---- | ---- | ---- | ---- |
|--avatar.bucket|aliyun oss bucket|string|""|
|--avatar.endpoint|aliyun oss endpoint|string|""|
|--avatar.accesskeyid|aliyun access key id|string|""|
|--avatar.accesskeysecret|aliyun access key secret|string|""|
|--avatar.cdn|cdn for aliyun oss|string|""|

### third party login

Now google and wechat login are supported.

|  Command line flag  |  Description  | Type | Default |
| ---- | ---- | ---- | ---- |
|--google_login.aud|audience of google login|string|""|
|--wechat_login.app_id|wechat app id|string|""|
|--wechat_login.secret|wechat secret|string|""|


### JWT

Expiration time of JWT can be configured.

|  Command line flag  |  Description  | Type | Default |
| ---- | ---- | ---- | ---- |
|--jwt.access_token_expire|expire time(s) of access token|int|600|
|--jwt.reset_password_token_expire|expire time(s) of reset password token|int|1200|
|--jwt.reset_password_result_token_expire|expire time(s) of reset password result token|int|300|
|--jwt.register_verify_token_expire|expire time(s) of reset password result token|int|300|

### view

View template location which Pluto used like register verfication mail or reset password page.

You can also replace the default pages with your own pages, see [docs/view.md](https://github.com/MuShare/pluto/blob/master/docs/view.md) for more information.

|  Command line flag  |  Description  | Type | Default |
| ---- | ---- | ---- | ---- |
|--view.path|path of html view location directory|string|./views|

## Use Config File

Pluto support config file in JSON, YAML, TOML format. All the flags above can be set through config file.

JSON config file example:
```
{
  "server": {
    "skip_register_verify_mail": true
  },
  "log": {
    "format": "logfmt",
    "level": "debug"
  },
  "database": {
    "host": "127.0.0.1",
    "port": "3306",
    "user": "root",
    "password": "root",
    "db": "pluto_server"
  },
  "rsa": {
    "path": "/etc/pluto",
    "name": "id_rsa_test"
  }
}
```
