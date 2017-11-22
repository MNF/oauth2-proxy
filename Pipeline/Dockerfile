ARG IMAGE_NAME
ARG IMAGE_VERSION

FROM $IMAGE_NAME:$IMAGE_VERSION as builder

# Build runtime image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh

WORKDIR /root/
RUN MKDIR /root/config
COPY --from=builder /go/src/oauth2_proxy/oauth2_proxy .
#ENTRYPOINT []
ENTRYPOINT ["/root/oauth2_proxy"]
