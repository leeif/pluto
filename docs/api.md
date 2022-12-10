# API Document
Pluto provides a set of HTTP resetful APIs.

- [API Document](#api-document)
  - [account api](#account-api)
    - [/v1/user/register](#v1userregister)
    - [/v1/user/register/verify/mail](#v1userregisterverifymail)
    - [/v1/user/login/account](#v1userloginaccount)
    - [/v1/user/login/google/mobile](#v1userlogingooglemobile)
    - [/v1/user/login/apple/mobile](#v1userloginapplemobile)
    - [/v1/user/login/wechat/mobile](#v1userloginwechatmobile)
    - [/v1/user/password/reset/mail](#v1userpasswordresetmail)
    - [/v1/user/binding](#v1userbinding)
    - [/v1/user/unbinding](#v1userunbinding)
    - [/v1/user/summary](#v1usersummary)
    - [/v1/user/count](#v1usercount)
    - [/v1/user/info:GET](#v1userinfoget)
    - [/v1/user/info:PUT](#v1userinfoput)
    - [/v1/user/delete](#v1userdelete)
  - [health check api](#health-check-api)
    - [/v1/healthcheck](#v1healthcheck)
  - [rbac api](#rbac-api)
    - [/v1/rbac/role/create](#v1rbacrolecreate)
    - [/v1/rbac/scope/create](#v1rbacscopecreate)
    - [/v1/rbac/role/scope](#v1rbacrolescope)
    - [/v1/rbac/role/scope/default](#v1rbacrolescopedefault)
    - [/v1/rbac/application/create](#v1rbacapplicationcreate)
    - [/v1/rbac/application/role/default](#v1rbacapplicationroledefault)
    - [/v1/rbac/application/list](#v1rbacapplicationlist)
    - [/v1/rbac/application/update-i18n-name](#v1rbacapplicationupdate-i18n-name)
    - [/v1/rbac/application/i18n-name](#v1rbacapplicationi18n-name)
    - [/v1/rbac/role/list](#v1rbacrolelist)
    - [/v1/rbac/scope/list](#v1rbacscopelist)
    - [/v1/rbac/user/application/role](#v1rbacuserapplicationrole)
  - [token api](#token-api)
    - [/v1/token/refresh](#v1tokenrefresh)
    - [/v1/token/publickey](#v1tokenpublickey)
    - [/v1/token/access/verify](#v1tokenaccessverify)
  - [oauth api](#oauth-api)
    - [/v1/oauth/tokens](#v1oauthtokens)
    - [/v1/oauth/client](#v1oauthclient)
    - [/v1/oauth/client/status](#v1oauthclientstatus)
  - [Api authorization](#api-authorization)
  - [Error Response](#error-response)


## account api

### /v1/user/register

Register user with email

* method: POST

* request:

```
{"mail":<string>, "name": <string>, "password":<string>, "app_id":<string>, "user_id":<string>}
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

### /v1/user/register/verify/mail

Send registration verification mail

* method: POST

* request:

```
{
    "mail": string,
    "app_id": string,
    "user_id": string
}
```

* response example:

```
{
  "status": "ok",
  "body": nil
}
```

### /v1/user/login/account

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

### /v1/user/login/google/mobile

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

### /v1/user/login/apple/mobile

Login with apple account for mobile app

* method: POST

* request:

```
{"code":<string>, "name":<string>, "device_id":<string>, "app_id":<string>}
```

Code is the token for verifying the register and getting the user's info like id and email

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

### /v1/user/login/wechat/mobile

Login with wechat account for mobile app

Offical docs [wechat](https://developers.weixin.qq.com/doc/oplatform/Mobile_App/WeChat_Login/Development_Guide.html).

* method: POST

* request:

```
{"code":<string>, "device_id":<string>, "app_id":<string>}
```

Code is the token for exchanging the access token of wechat

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

### /v1/user/password/reset/mail

Send password reset mail

* method: POST

* request:

```
{
    "mail": <string>,
    "app_id": <string>
}
```

* response example:

```
{
  "status": "ok",
  "body": nil
}
```

### /v1/user/binding

Bind mail, google, wechat, apple account

* method: POST

* request:

```
{"type":"mail|google|apple|wechat", "id_token":"<string>"(google binding), "mail":"<string>"(mail binding), "code":"<string>"(wechat, apple binding)}
```

### /v1/user/unbinding

Unbind mail, google, wechat, apple account

* method: POST

* request:

```
{"type":"mail|google|apple|wechat"}
```

```
{"mail":<string>}
```

### /v1/user/summary

Search the user using user name

* Access token with `pluto.admin` scope needs

* method: GET

### /v1/user/count

Get the count of the total users

* Access token with `pluto.admin` scope needs

* method: GET

### /v1/user/info:GET

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
    "verified": true,
    "user_id": xxasdfw,
    "user_id_updated": false
  }
}
```

### /v1/user/info:PUT

Update user info

* Access token needs

* method: PUT

* request:

```
{"name":<string>, "gender":<string>, "avatar":<string>, "user_id":<string>}
```

* response example:

```
{
  "status": "ok",
  "body": nil
}
```

### /v1/user/delete

Delete user

* Access token needs

* method: POST

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

Access token with `pluto.admin` scope should be provided

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
### /v1/rbac/application/update-i18n-name
update application i18n name

* method: POST

* request:

```
{
    "app_id": int,
    "i18n_names":[
        {
            "tag": string,
            "i18n_name": string
        }
    ]
}
```
### /v1/rbac/application/i18n-name
get application i18n names

* method: GET

* request query parameters:
```
"app_id": int
```
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

* Access token with `pluto.admin` scope needs

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
