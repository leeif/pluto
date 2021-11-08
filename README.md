# Pluto

[![Build Status](https://travis-ci.org/mushare/pluto.svg?branch=master)](https://travis-ci.org/mushare/pluto)
[![Go Report Card](https://goreportcard.com/badge/github.com/MuShare/pluto)](https://goreportcard.com/report/github.com/MuShare/pluto)
[![Gitter](https://badges.gitter.im/pluto-discuss/community.svg)](https://gitter.im/pluto-discuss/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

Pluto is a JWT based authorization/authentication service. Besides providing a basic user registration and login feature, Pluto also provides a RBAC management to control the user's permission. Pluto implements the OAuth2 specified APIs for authorization.

## Setup

Go versio: 1.16

```bash
# install sqlboiler
$ go install -v github.com/volatiletech/sqlboiler@v3.6.0
$ go install -v github.com/volatiletech/sqlboiler/drivers/sqlboiler-mysql@v3.6.0
$ git clone ...
# start mysql and create `pluto` database
$ make migrate-binary-build
# run migrations
$ ./bin/pluto-migrate
# build server
$ make server-binary-install
# start server
$ pluto-server
```

## Main Features

* User registration / login
* Oauth2 APIs
* JWT-based authorization
* Role-based access control (RBAC)
* Admin page [link](https://github.com/MuShare/pluto-admin)

## Getting started

The [Installation doc](https://github.com/MuShare/pluto/blob/master/docs/installation.md) have a guide on how to setup the Pluto server via Docker images, Kubernetes or from source.

### Documents

All documents can be found in [/docs](https://github.com/MuShare/pluto/blob/master/docs)

Here are some helpful documents for reading.

* [API Document](https://github.com/MuShare/pluto/blob/master/docs/api.md)
* [Oauth2](https://github.com/MuShare/pluto/blob/master/docs/oauth.md)
* [Configuration](https://github.com/MuShare/pluto/blob/master/docs/configuration.md)
* [Replace Views](https://github.com/MuShare/pluto/blob/master/docs/view.md) is a guide for replacing the default html pages with your own custom files
* [JWT Token](https://github.com/MuShare/pluto/blob/master/docs/jwt.md) gives an introduction of the JWT design in Pluto.
* [WeChat Login](https://github.com/MuShare/pluto/blob/master/docs/wechat.md) gives an introduction of signing in with WeChat QRCode.

## Docker image

https://hub.docker.com/repository/docker/mushare/pluto

## Contribute

Feel free to fire an issue or send a pull request.

## License

MIT License, see [LICENSE](https://github.com/MuShare/pluto/blob/master/LICENSE)
