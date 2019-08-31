FROM golang:1.12

ENV ConfigFile /etc/pluto/config.json

ADD . /go/src/pluto

WORKDIR /go/src/pluto

RUN  export GO111MODULE=on && go build -o pluto-server cmd/pluto-server/main.go

CMD ./pluto-server --config.file=$ConfigFile