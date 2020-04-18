# API Document
Pluto provides a set of HTTP resetful APIs.

* [`/v1/register`](#v1register)
* [`/v1/register/verify/mail`](#v1registerverifymail)
* [`/v1/login/account`](#v1loginaccount)
* [`/v1/login/google/mobile`](#v1logingooglemobile)
* [`/v1/login/apple/mobile`](#v1loginapplemobile)
* [`/v1/login/wechat/mobile`](#v1loginwechatmobile)
* [`/v1/password/reset/mail`](#v1passwordresetmail)
* [`/v1/healthcheck`](#v1healthcheck)
* [`/v1/rbac/role/create`](#v1rbacrolecreate)
* [`/v1/rbac/scope/create`](#v1rbacscopecreate)
* [`/v1/rbac/role/scope`](#v1rbacrolescope)
* [`/v1/rbac/role/scope/default`](#v1rbacrolescopedefault)
* [`/v1/rbac/application/create`](#v1rbacapplicationcreate)
* [`/v1/rbac/application/role/default`](#v1rbacapplicationroledefault)
* [`/v1/rbac/application/list`](#v1rbacapplicationlist)
* [`/v1/rbac/role/list`](#v1rbacrolelist)
* [`/v1/rbac/scope/list`](#v1rbacscopelist)
* [`/v1/rbac/user/application/role`](#v1rbacuserapplicationrole)
* [`/v1/user/search`](#v1usersearch)
* [`/v1/user/count`](#v1usercount)
* [`/v1/user/info/{userID}`](#v1userinfo{userID})
* [`/v1/user/info/{userID}`](#v1userinfo{userID})
* [`/v1/token/refresh`](#v1tokenrefresh)
* [`/v1/token/publickey`](#v1tokenpublickey)
* [`/v1/token/access/verify`](#v1tokenaccessverify)
* [`/v1/oauth/tokens`](#v1oauthtokens)
* [`/v1/oauth/client`](#v1oauthclient)
* [`/v1/oauth/client/status`](#v1oauthclientstatus)


## account api

### /v1/register

Register user with email

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
    "mail": "geeklyf@hotmail.com",
    "verified": false
  }
}
 ```

### /v1/register/verify/mail

Send registration verification mail

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

### /v1/login/account

Login with email or username

* method: POST

* request:

```
{"account":<string>, "password":<string>, "device_id":<string|optional>, "app_id":<string>}
```

* response example:

```
{
  "status": "ok",
  "body": {
    "refresh_token": "",
    "access_token": "",
    "type": "Bearer"
  }
}
```

### /v1/login/google/mobile

Login with google account for mobile app

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
    "refresh_token": "",
    "access_token": "",
    "type": "Bearer"
  }
}
```

### /v1/login/apple/mobile

Login with apple account for mobile app

* method: POST

* request:

```
{"code":<string>, "name":<string>, "device_id":<string>, "app_id":<string>}
```

Code is the token for verifying the register and getting the user's info like id and email.

* response example:

```
{
  "status": "ok",
  "body": {
    "refresh_token": "",
    "access_token": "",
    "type": "Bearer"
  }
}
```

### /v1/login/wechat/mobile

Login with wechat account for mobile app

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
    "refresh_token": "",
    "access_token": "",
    "type": "Bearer"
  }
}
```

### /v1/password/reset/mail

Send password reset mail

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

## health check api

### /v1/healthcheck

Health check

* method: Get

* response example:

```
{
  "status": "ok",
  "body": {
    "version": ""
  }
}
```

## rbac api

Access token with `pluto.admin` scope should be provided.

### /v1/rbac/role/create

Create role

* method: POST

* request:

```
{"app_id":<string>, "name":<string>}
```

* response example:

### /v1/rbac/scope/create

Create scope

* method: POST

* request:

```
{"app_id":<string>, "name":<string>}
```

* response example:

### /v1/rbac/role/scope

Update scopes of the role

* method: POST

* request:

```
{"role_id":<int>, "scopes":[array...]}
```

* response example:

### /v1/rbac/role/scope/default

Set the default scope of the role

* method: POST

* request:

```
{"role_id":<int>, "scope_id":<int>}
```

* response example:

### /v1/rbac/application/create

Create application

* method: POST

* request:

```
{"name":<string>}
```

* response example:

### /v1/rbac/application/role/default

Set the default role of the application

* method: POST

* request:

```
{"app_id":<string>, "role_id":<string>}
```

* response example:

### /v1/rbac/application/list

List all the applications

* method: GET

* request:

```
{"mail":<string>}
```

* response example:

### /v1/rbac/role/list

List all the roles in the application

* method: GET

* request:

```
{"app_id":<string>}
```

* response example:

### /v1/rbac/scope/list

List all the scopes in the application

* method: GET

* request:

```
{"app_id":<string>}
```

* response example:

### /v1/rbac/user/application/role

Set the role of a user in application

* method: POST
* request:

```
{"user_id":<string>, "app_id":<int>, "role_id":<int>}
```

* response example:

## user api

### /v1/user/search

Search the user using name or mail

* Access token with `pluto.admin` scope needs.

* method: GET

* request:

```
{"account":<string>}
```

### /v1/user/count

Get the count of the total users

* Access token with `pluto.admin` scope needs.

* method: GET

### /v1/user/info/{userID}

Get user info

* Access token needs

* method: GET

* response example:

```
{
  "status": "ok",
  "body": {
    "avatar": "",
    "created_at": 1586925495,
    "login_type": "mail",
    "mail": "geeklyf92610@gmail.com",
    "name": "yifan.li",
    "roles": "admin",
    "sub": 1,
    "updated_at": 1586925495,
    "verified": true
  }
}
```

### /v1/user/info/{userID}

Update user info

* Access token needs

* method: PATCH

* request:

```
{"name":<string>, "gender":<string>, "avatar":<string>}
```

* response example:

```
{
  "status": "ok",
  "body": nil
}
```



## token api

### /v1/token/refresh

Refresh access token

* method: POST
* request:

```
{"refresh_token":<string>, "app_id":<string>, "scopes":<string|optional>}
```

* response example:

```
{
  "status": "ok",
  "body": {
    "refresh_token": "",
    "access_token": "",
    "type": "Bearer"
  }
}
 ```

### /v1/token/publickey

Get the rsa public key

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

### /v1/token/access/verify

Verify access token

* method: Get

* request:

```
{"token":<string>}
```

* response example:

```
{
  "status": "ok",
  "body": {
    "type": "access",
    "iat": 1587041608,
    "exp": 1587045208,
    "sub": 3,
    "iss": "pluto",
    "scopes": [
        ""
    ]
  }
}
```

## oauth api

### /v1/oauth/tokens

Request access token

### /v1/oauth/client

Create client

* Access token needs.

* method: POST

* request:

```
{"key":<string>, "secret":<string>}
```

### /v1/oauth/client/status

Set the client status to approved

* Access token with `pluto.admin` scope needs.

* method: PUT

* request:

```
{"key":<string>, "status":<approved|denied>}
```

## Api authorization

Some APIs above should be authorized with an access token.

Pluto use the Bearer authorize which you can set in the Authorization header.

* header: Authorization: Bearer \<access token\>

## Error Response

The error responsed by Pluto is in the json format as below:

```
{
  "status": "error",
  "error": {
    "code": <error code>,
    "message": "<error message>"
  }
}
```
