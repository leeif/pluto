# Installation

## Setup Database for Pluto

Pluto provides a database migration tool to create database tables.

```
git clone https://github.com/MuShare/pluto.git
cd pluto/
make migrate-binary-build

# set the database config in config.yaml
./bin/pluto-migrate --config.file=config.yaml
```

## Installation from source

Requirement: Go (> 1.13)

Download the source file.
```
git clone https://github.com/MuShare/pluto.git
cd pluto/
make server-binary-build
./bin/pluto-server --config.file=config.yaml
```

## Run Pluto with Docker image
```
docker run -v ${PWD}/config.json:/etc/pluto/config.json --name pluto -d mushare/pluto:latest
docker logs pluto
```
The default config file location in Pluto container is /etc/pluto/config.json.
Here we mount a local config file into container. About config, please see here [Configuration](https://github.com/MuShare/pluto/blob/master/docs/configuration.md) for more information.

Pluto also support several other formats of config file like YAML, TOML.
To use these kinds , you need to pass an env variable like this,
```
docker run --env ConfigFile=/etc/pluto/config.yaml -v ${PWD}/config.yaml:/etc/pluto/config.yaml --name pluto -d mushare/pluto:latest
```

## Run Pluto in Kubernetes
create a configmap from a json config file.
```
kubectl create configmap pluto-config --from-file=config.json
```
Deployment
```
cat << EOS
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pluto-deployment
spec:
  replicas: 2
  selector:
    matchLabels:
      app: pluto
  template:
    metadata:
      labels:
        app: pluto
    spec:
      containers:
      - name: pluto
        image: mushare/pluto:latest
        volumeMounts:
        - name: config-volume
          mountPath: /etc/pluto/config.json
          subPath: config.json
        ports:
        - containerPort: 8010
      volumes:
      - name: config-volume
        configMap:
          name: pluto-config
EOS | kubectl apply -f -
```

## Setup Your First Application

Pluto is designed as an application-based auth system.

You can start creating your first application through pluto admin page.

### Start the admin web page

### Login the admin web page

### Create the Application

### Create the Role and Scopes
