FROM golang:1.9 AS builder
WORKDIR /go/src/oauth2_proxy
COPY . /go/src/oauth2_proxy/
RUN go get -x -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

