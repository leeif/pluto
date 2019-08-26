# Pluto Server

[![Build Status](https://travis-ci.org/MuShare/pluto.svg?branch=master)](https://travis-ci.org/MuShare/pluto)
[![Go Report Card](https://goreportcard.com/badge/github.com/leeif/pluto)](https://goreportcard.com/report/github.com/leeif/pluto)

# API Document

## 1. User
(1) /api/user/register

 * method: POST
 * request: 
 ```
 {"mail":<string>, "name": <string>, "password":<string>}
 ```
 * response:
 ```
 {
  "status": "ok",
  "body": {
    "mail": "geeklyf@hotmail.com"
  }
}
 ```

(2) /api/user/register/verify/mail

resend register verify mail
 * method: POST
 * request:
 ```
 {"mail":<string>}
 ```
 * response:
 ```
{
  "status": "ok",
  "body": nil
}
 ```

(3) /api/user/login

* method: POST
* request:
```
{"mail":<string>, "password":<string>, "device_id":<string>, "app_id":<string>}
```
* response:
```
{
  "status": "ok",
  "body": {
    "jwt": "",
    "refresh_token": ""
  }
}
```

(4) /api/user/password/reset/mail

send password reset form to mail
* method: POST
* request:
```
{"mail":<string>}
```
* response:
```
{
  "status": "ok",
  "body": nil
}
```

(5) /api/user/info/\<token\>

Get User info, token needs a base64 encoded string of jwt token
* method: GET
* response:
```
{
  "status": "ok",
  "body": {
    "ID": 1,
    "CreatedAt": "2019-08-26T13:35:32+09:00",
    "UpdatedAt": "2019-08-26T13:35:39+09:00",
    "DeletedAt": null,
    "Mail": "geeklyf@hotmail.com",
    "Name": "yifan.li",
    "Gender": null,
    "Birthday": null,
    "Avatar": "https://pluto-staging.oss-cn-hongkong.aliyuncs.com/avatar/db226f6bb95c1c5e.png"
  }
}
```

## 2. Auth
(1) /api/auth/publickey

get the public key from auth server
 * method: Get
 * response:
 ```
{
  "status": "ok",
  "body": {
    "public_key": ""
  }
}
 ```
(2) /api/auth/refresh

get a new jwt access token
* method: POST
* request:
```
{"refresh_token":<string>, "user_id":<int>, "device_id":<string>, "app_id":<string>}
```
* response:
 ```
 {
  "status": "ok",
  "body": {
    "jwt": ""
  }
}
 ```

# Access JWT Token Example

expire: 3600s
```
{
 type: "access",
 alg: "rsa"
}.
{
 create_time: 1566396559,
 expire_time: 1566400159,
 userId: 1,
 deviceId: "xxxx",
 appId: "xxxxx"
}.
[signature]
```