# view
Pluto is using the golang [template](https://golang.org/pkg/text/template/) engine to parse the template HTML files.
Pluto default uses the templates in [views](https://github.com/MuShare/pluto/blob/master/views).


## template/
All the file in this directory will be loaded in every template parsing, you can put header, footer templates in here.

## error.html
The error.html will be return if there is an server internal error or JWT authentication error.

## password_reset_mail.html
Password reset mail HTML template.

Data
```
BaseURL <string> : base url of request
Token <string> : JWT token
```

## password_reset.html
Password reset page.
Data
```
Error <pluto_error>
Token <string>
```

## password_reset_result.html
Password reset result page.

Data
```
Error <pluto_error>
```

## register_verify_mail.html
Register verification mail.

Data
```
BaseURL <string> : base url of request
Token <string> : JWT token
```

## register_verify_result.html
Register verification result page.

Data
```
Error <pluto_error>
```

# Replace the default views
Docker
```
docker run -v ./config.json:/etc/pluto/config.json -v ./views:/views --name pluto -d leeif/pluto:latest
```
Replace the default views in container, remember the name of your files should be the same as above except the files under the [views/template/](https://github.com/MuShare/pluto/blob/master/views/template/).