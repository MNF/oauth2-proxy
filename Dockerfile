FROM containerregistrydev.azurecr.io/library/golang:1.15-alpine AS builder

# installing bash
RUN apk update && apk upgrade && \
    apk add --no-cache bash curl dep git

WORKDIR /go/src/oauth2_proxy
COPY . /go/src/oauth2_proxy/
RUN ["/bin/bash", "-c", "bash < <(curl -s -S -L https://raw.githubusercontent.com/golang/dep/v0.5.0/install.sh)"]
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .


# Build runtime image
FROM containerregistrydev.azurecr.io/library/alpine:3.10
RUN apk --no-cache add ca-certificates
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /root/
RUN mkdir /root/config
COPY --from=builder /go/src/oauth2_proxy/oauth2_proxy .

ENTRYPOINT ["/root/oauth2_proxy"]
