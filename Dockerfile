FROM golang:1.9 AS builder
WORKDIR /go/src/oauth2_proxy
COPY . /go/src/oauth2_proxy/
RUN ["/bin/bash", "-c", "bash < <(curl -s -S -L https://raw.githubusercontent.com/golang/dep/v0.5.0/install.sh)"]
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .
