# syntax = docker/dockerfile:1-experimental
FROM --platform=${BUILDPLATFORM} golang:alpine AS gobase
RUN apk add --no-cache \
        ca-certificates \
        gcc \
        musl-dev \
    && true
RUN go install golang.org/x/lint/golint@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
WORKDIR /src



FROM --platform=${BUILDPLATFORM} gobase AS gobuilder
ENV CGO_ENABLED=0
COPY go.mod go.sum .
RUN --mount=type=cache,target=/root/.cache/go-build go mod download
ARG BUILD_VERSION=master
ARG BUILD_COMMIT=unknown
ARG BUILD_DATE=now
ARG TARGETOS
ARG TARGETARCH
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build golint -set_exit_status ./...
RUN --mount=type=cache,target=/root/.cache/go-build go vet -v ./...
RUN --mount=type=cache,target=/root/.cache/go-build go test -v -short ./...
