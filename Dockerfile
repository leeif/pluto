FROM golang:1.12 as build

ARG VERSION

ADD . /go/src/pluto

WORKDIR /go/src/pluto

RUN  export GO111MODULE=on && go build -ldflags="-X 'main.VERSION=${VERSION}'" -o pluto-server cmd/pluto-server/main.go

FROM ubuntu:18.04

RUN apt-get update && apt-get install -y wget

ENV DOCKERIZE_VERSION v0.6.1
RUN wget https://github.com/jwilder/dockerize/releases/download/$DOCKERIZE_VERSION/dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && tar -C /usr/local/bin -xzvf dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz \
    && rm dockerize-linux-amd64-$DOCKERIZE_VERSION.tar.gz

ENV ConfigFile /etc/pluto/config.json

COPY --from=build /go/src/pluto/pluto-server /usr/bin/

COPY --from=build /go/src/pluto/views views/

CMD /usr/bin/pluto-server --config.file=$ConfigFile