# JWT

Pluto assign a JWT for each user and you can use this JWT to authenticate users in your own apps.

## JWT Token Example

```
{
 type: "jwt",
 alg: "rsa"
}.
{
  "type": "access",
  "iat": 1576230404,
  "exp": 1576234004,
  "sub": 1, // user id
  "iss": "pluto", // application name
  "scopes": [
    "pluto.admin"
  ],
  "login_type": "mail"
}.
[signature]
```

Here is the example of an access JWT token. You can customize the expiration time through [Configuration](https://github.com/MuShare/pluto/blob/master/README.md).

## Signature

### rsa algorithm

The third part of the JWT token is a signature signed with the rsa private key which provided through [rsa config](https://github.com/MuShare/pluto/blob/master/docs/configuration.md#rsa)

The sign text is a concat of head and payload of a JWT.
```
sign(string(head)+string(payload), <private key>)
```

Verify a signature need the rsa public key, which you can get through the [public key api](https://github.com/MuShare/pluto/blob/master/docs/api.md#apiauthpublickey)

```
verify(string(head)+string(payload), <public key>)
```


## Type of JWT in Pluto

In Pluto, we have the following types of JWT token.

* ACCESS token will be responsed after user login or oauth authorization.

*	REGISTERVERIFY token is to verify the register mail.

*	PASSWORDRESET token is used to reset the password.

Accept the ACCESS token, all other tokens are used internally in Pluto.
Since Pluto is based on JWT, we use JWT to access all the exposed HTML pages in Pluto.
