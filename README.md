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

(4) /api/user/login/google

* method: POST
* request:
```
{"id_token":<string>, "device_id":<string>, "app_id":<string>}
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

(5) /api/user/password/reset/mail

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

(6) /api/user/info/me

Get User info, token needs a base64 encoded string of jwt token
* method: GET
* header: Authorization: jwt base64(token)
* response:
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

## 3. Health Check
/healthcheck
* method: Get

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