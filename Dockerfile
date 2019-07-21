# This is a multi-stage Dockerfile and requires >= Docker 17.05
# https://docs.docker.com/engine/userguide/eng-image/multistage-build/
############################
# STEP 1 build executable binary
############################
FROM golang:alpine as builder

# Export app path as environment variable
ENV APP_PATH $GOPATH/src/github.com/bsmmoon/go-proxy
WORKDIR $APP_PATH
COPY . $APP_PATH

# Install tools for building purpose
RUN apk update \
    && apk add --virtual build-dependencies \
        build-base \
        gcc \
        wget \
        git \
        g++ \
        dep \
    && apk add \
        bash

RUN make install
RUN make build

############################
# STEP 2 build a small image
############################
# Build from the most lightweight image
FROM alpine:3.7
LABEL project="go-proxy"
LABEL image="build"

# Copy generated binary from builder image
COPY --from=builder /go/bin/proxy /go/bin/proxy

CMD /go/bin/proxy
