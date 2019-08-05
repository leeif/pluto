# Pluto Server

[![Build Status](https://travis-ci.org/MuShare/pluto.svg?branch=master)](https://travis-ci.org/MuShare/pluto)

# API Document

## 1. User
(1) /api/user/register

 * method: POST
 * body: {"mail":" ", "name": " ", "password":" "}

(2) /api/user/login

* method: POST
* body: {"mail":" ", "password":" ", "device_id":" ", "app_id":" "}

## 2. Auth
(1) /api/auth/publickey

 * method: Get

(2) /api/auth/refresh

* method: POST
* body: {"refresh_token":" ", "user_id":" ", "device_id":" ", "app_id":" "}