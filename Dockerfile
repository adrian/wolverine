#
# Build Stage
#
FROM golang:1.18.0-alpine3.15 as build

ENV APP_NAME wolverine
ENV CMD_PATH cmd/main.go

ARG VERSION=dev

COPY . $GOPATH/src/$APP_NAME
WORKDIR $GOPATH/src/$APP_NAME
 
RUN go build -ldflags="-X main.BuildVersion=${VERSION}" -v -o /$APP_NAME $GOPATH/src/$APP_NAME/$CMD_PATH

#
# Run Stage
#
FROM alpine:3.15
 
ARG METRICS_PORT=2112

ENV APP_NAME wolverine

COPY --from=build /$APP_NAME .
 
EXPOSE $METRICS_PORT
 
CMD ./$APP_NAME
