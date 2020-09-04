# Pluto

[![Build Status](https://travis-ci.org/leeif/pluto.svg?branch=master)](https://travis-ci.org/leeif/pluto)
[![Go Report Card](https://goreportcard.com/badge/github.com/MuShare/pluto)](https://goreportcard.com/report/github.com/MuShare/pluto)
[![Gitter](https://badges.gitter.im/pluto-discuss/community.svg)](https://gitter.im/pluto-discuss/community?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

Pluto is a JWT based authorization/authentication service. Besides providing a basic user registration and login feature, Pluto also provides a RBAC management to control the user's permission. Pluto implements the OAuth2 specified APIs for authorization.

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

## Contribute

Feel free to fire an issue or send a pull request.

## License

MIT License, see [LICENSE](https://github.com/MuShare/pluto/blob/master/LICENSE)
