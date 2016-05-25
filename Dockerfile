FROM golang:1.6.2-alpine
MAINTAINER Phil Estes <estesp@gmail.com>

ENV GOPROJ /go/src/github.com/estesp/elastistack

ENV PATH /go/bin:/usr/local/go/bin:$PATH
ENV GOPATH /go:$GOPROJ/Godeps/_workspace

WORKDIR $GOPROJ
COPY . $GOPROJ

RUN go build -o /go/bin/elastistack .

ENTRYPOINT [ "elastistack" ]
