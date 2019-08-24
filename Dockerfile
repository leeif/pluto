FROM golang:1.12

ENV ConfigFile /etc/pluto/config.json

RUN apt-get update && apt-get install -y go-dep

ADD . /go/src/pluto

WORKDIR /go/src/pluto

RUN dep ensure && go build -o pluto-server cmd/pluto-server/main.go

CMD ./pluto-server --config.file=$ConfigFile