FROM golang:1.16-buster AS builder
ARG VERSION

# Copy sources
WORKDIR $GOPATH/src/github.com/oauth2-proxy/oauth2-proxy

# Fetch dependencies
COPY go.mod go.sum ./
RUN GO111MODULE=on go mod download

# Now pull in our code
COPY . .

# Build binary and make sure there is at least an empty key file.
#  This is useful for GCP App Engine custom runtime builds, because
#  you cannot use multiline variables in their app.yaml, so you have to
#  build the key into the container and then tell it where it is
#  by setting OAUTH2_PROXY_JWT_KEY_FILE=/etc/ssl/private/jwt_signing_key.pem
#  in app.yaml instead.
RUN VERSION=${VERSION} make build && touch jwt_signing_key.pem

# Copy binary to alpine
FROM alpine:3.13
COPY nsswitch.conf /etc/nsswitch.conf
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/oauth2-proxy/oauth2-proxy/oauth2-proxy /bin/oauth2-proxy
COPY --from=builder /go/src/github.com/oauth2-proxy/oauth2-proxy/jwt_signing_key.pem /etc/ssl/private/jwt_signing_key.pem

#for local tests (TODO: parameterise if will cause problems on dev/prod deployments)
RUN mkdir /etc/config/
COPY pipeline/oauth2_proxy.local.cfg /etc/config/oauth2_proxy.cfg 
RUN ls -l /etc/config/

#RUN export HTTP_PROXY=127.0.0.1:8888
#RUN export HTTPS_PROXY=127.0.0.1:8888
#https://stackoverflow.com/questions/54218632/how-to-use-local-proxy-settings-in-docker-compose/57714235#57714235
#ENV http_proxy http://127.0.0.1:8888
#[2021/03/30 08:33:15] [oauthproxy.go:737] Error redeeming code during OAuth2 callback: token exchange failed: Post "https://webjetb2cdev.b2clogin.com/webjetb2cdev.onmicrosoft.com/B2C_1A_signin_TSA_Local_NZ/oauth2/v2.0/token": proxyconnect tcp: dial tcp 127.0.0.1:8888: connect: connection refused
#ENV https_proxy http://127.0.0.1:8888

USER 2000:2000
#suggestion from https://github.com/microsoft/vscode-remote-release/issues/174#issuecomment-489917484
#ENV HOME /home/node

ENTRYPOINT ["/bin/oauth2-proxy"]
