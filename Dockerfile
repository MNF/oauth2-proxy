FROM golang:1.9 AS builder
WORKDIR /go/src/oauth2_proxy
COPY . /go/src/oauth2_proxy/
RUN go get -x -d -v
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/oauth2_proxy/oauth2_proxy .
ENTRYPOINT ["/root/oauth2_proxy"]
