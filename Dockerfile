##################
# Building image #
##################
FROM golang:1.10-alpine AS build-env

RUN mkdir -p /go/src/github.com/AirGateway/agw-proxy
WORKDIR /go/src/github.com/AirGateway/agw-proxy
ADD . /go/src/github.com/AirGateway/agw-proxy

RUN apk update && \
  apk add ca-certificates curl git && \
  rm -rf /var/cache/apk/*

RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh && \
  dep ensure && \
  go build -o agw-proxy

###############
# Final image #
###############
FROM alpine

RUN apk update && \
  apk add ca-certificates curl && \
  rm -rf /var/cache/apk/*

WORKDIR /app
COPY --from=build-env /go/src/github.com/AirGateway/agw-proxy/agw-proxy /app/agw-proxy

ENTRYPOINT ./agw-proxy
