# Installation

## Run Pluto with Docker image
```
docker run -v ./config.json:/etc/pluto/config.json --name pluto -d leeif/pluto:latest
docker logs pluto
```
The default config file location in Pluto container is /etc/pluto/config.json.
Here we mount a local config file into container. About config, please see here [Configuration](https://github.com/MuShare/pluto/blob/master/docs/configuration.md) for more information.

Pluto also support several other formats of config file like YAML, TOML.
To use these kinds , you need to pass an env variable like this,
```
docker run --env ConfigFile=/etc/pluto/config.yaml -v ./config.yaml:/etc/pluto/config.yaml --name pluto -d leeif/pluto
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
        image: leeif/pluto:latest
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

## Installation from source

Requirement: Go (> 1.12)

Download the source file.
```
git clone https://github.com/MuShare/pluto.git
cd pluto/
```
Build the binary
```
make binary-build
```
The binary will be build into ./bin
```
./bin/pluto-server --config.file=config.yaml
```
