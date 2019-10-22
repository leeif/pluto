# JWT

Pluto assign a JWT for each user and you can use this JWT to authenticate users in your own apps.

## JWT Token Example

```
{
 type: "jwt",
 alg: "rsa"
}.
{
 create_time: 1566396559,
 expire_time: 1566400159,
 type: "access",
 userId: 1,
 deviceId: "xxxx",
 appId: "xxxxx"
}.
[signature]
```
Here is the example of an access JWT token. You can customize the expiration time through [Configuration](https://github.com/mushare/kiper/blob/master/README.md).

## Type of JWT in Pluto
In Pluto, we have the following types of JWT token.

* ACCESS token will be response after user login.
*	REGISTERVERIFY token is to verify the register mail.
*	PASSWORDRESET token is used to reset the password.
*	PASSWORDRESETRESULT token is used to access the result page of password reset.

Accept the ACCESS token, all other tokens are used internally in Pluto.
Since Pluto is based on JWT, we use JWT to access all the exposed HTML pages in Pluto.
