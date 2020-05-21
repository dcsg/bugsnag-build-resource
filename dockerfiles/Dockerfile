FROM golang:1 AS builder
COPY . /go/src/github.com/dcsg/bugsnag-build-resource
ENV CGO_ENABLED 0
RUN go get -d ./...
RUN go build -o /assets/out github.com/dcsg/bugsnag-build-resource/cmd/out
RUN go build -o /assets/in github.com/dcsg/bugsnag-build-resource/cmd/in
RUN go build -o /assets/check github.com/dcsg/bugsnag-build-resource/cmd/check

FROM alpine:edge AS resource
RUN apk add --update bash ca-certificates
COPY --from=builder /assets /opt/resource
RUN chmod +x /opt/resource/*

FROM resource
