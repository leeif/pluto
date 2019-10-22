# API Document
Pluto now provides a set of HTTP resetful APIs.

They are as follows:

* [`/api/user/register`](#/api/user/register)
* [`/api/user/register/verify/mail`](#/api/user/register/verify/mail)
* [`/api/user/login`](#/api/user/login)
* [`/api/user/login/google/mobile`](#/api/user/login/google/mobile)
* [`/api/user/login/wechat/mobile`](#/api/user/login/wechat/mobile)
* [`/api/user/password/reset/mail`](#/api/user/password/reset/mail)
* [`/api/user/info/me`](#/api/user/info/me)
* [`/api/auth/publickey`](#/api/auth/publickey)
* [`/api/auth/refresh`](#/api/auth/refresh)
* [`/healthcheck`](#/healthcheck)

## User

### /api/user/register

Register pluto with personal mail.

 * method: POST
 * request: 
 ```
 {"mail":<string>, "name": <string>, "password":<string>}
 ```
 * response example:
 ```
 {
  "status": "ok",
  "body": {
    "mail": "geeklyf@hotmail.com"
  }
}
 ```

### /api/user/register/verify/mail

resend register verify mail
 * method: POST
 * request:
 ```
 {"mail":<string>}
 ```
 * response example:
 ```
{
  "status": "ok",
  "body": nil
}
 ```

### /api/user/login

Login with personal mail.

* method: POST
* request:
```
{"mail":<string>, "password":<string>, "device_id":<string>, "app_id":<string>}
```
* response example:
```
{
  "status": "ok",
  "body": {
    "jwt": "",
    "refresh_token": ""
  }
}
```

### /api/user/login/google/mobile

Login with google account for mobile APPs.

Offical docs [iOS](https://developers.google.com/identity/sign-in/ios/backend-auth), 
[Android](https://developers.google.com/identity/sign-in/android/backend-auth).

* method: POST
* request:
```
{"id_token":<string>, "device_id":<string>, "app_id":<string>}
```
* response example:
```
{
  "status": "ok",
  "body": {
    "jwt": "",
    "refresh_token": ""
  }
}
```

### /api/user/login/wechat/mobile

Login with wechat account for mobile APPs.

Offical docs [wechat](https://developers.weixin.qq.com/doc/oplatform/Mobile_App/WeChat_Login/Development_Guide.html).

* method: POST
* request:
```
{"code":<string>, "device_id":<string>, "app_id":<string>}
```
Code is the token for exchanging the access token of wechat.
* response example:
```
{
  "status": "ok",
  "body": {
    "jwt": "",
    "refresh_token": ""
  }
}
```

### /api/user/password/reset/mail

Send password reset form to mail.

* method: POST
* request:
```
{"mail":<string>}
```
* response example:
```
{
  "status": "ok",
  "body": nil
}
```

### /api/user/info/me

Get User info which requires a base64 encoded string of jwt token.

* method: GET
* header: Authorization: jwt base64(token)
* response example:
```
{
  "status": "ok",
  "body": {
    "id": 2,
    "create_at": 1567236370,
    "updated_at": 1567236410,
    "delete_at": 0,
    "mail": "geeklyf92610@gmail.com",
    "name": "yifan.li",
    "gender": null,
    "birthday": 0,
    "avatar": "xxxx"
  }
}
```

## Auth

### /api/auth/publickey

Get the public key from auth server.
 * method: Get
 * response example:
 ```
{
  "status": "ok",
  "body": {
    "public_key": ""
  }
}
 ```

### /api/auth/refresh

Get a new jwt access token.
* method: POST
* request example:
```
{"refresh_token":<string>, "user_id":<int>, "device_id":<string>, "app_id":<string>}
```
* response example:
 ```
 {
  "status": "ok",
  "body": {
    "jwt": ""
  }
}
 ```

## Health Check

### /healthcheck
* method: Get

* response example:
 ```
 {
  "status": "ok",
  "body": null
}
