FROM golang:1.12 as build

ARG VERSION

ADD . /go/src/pluto

WORKDIR /go/src/pluto

RUN  export GO111MODULE=on && go build -ldflags="-X 'main.VERSION=${VERSION}'" -o pluto-server cmd/pluto-server/main.go

FROM ubuntu:18.04

ENV ConfigFile /etc/pluto/config.json

COPY --from=build /go/src/pluto/pluto-server /usr/bin/

CMD /usr/bin/pluto-server --config.file=$ConfigFile