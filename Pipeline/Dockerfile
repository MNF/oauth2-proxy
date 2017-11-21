ARG IMAGE_NAME
ARG IMAGE_VERSION

FROM $IMAGE_NAME:$IMAGE_VERSION as builder

# Build runtime image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/oauth2_proxy/oauth2_proxy .
ENTRYPOINT ["/root/oauth2_proxy"]
