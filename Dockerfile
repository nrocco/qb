# syntax = docker/dockerfile:1-experimental
FROM --platform=${BUILDPLATFORM} golang:alpine AS base
RUN apk add --no-cache \
        ca-certificates \
        gcc \
        musl-dev \
    && true
RUN env GO111MODULE=on go get -u \
        golang.org/x/lint/golint \
        golang.org/x/tools/cmd/goimports \
        && true
WORKDIR /app

FROM base AS builder
ENV CGO_ENABLED=1
COPY go.mod go.sum .
RUN go mod download
ARG BUILD_VERSION=master
ARG BUILD_COMMIT=unknown
ARG BUILD_DATE=now
ARG TARGETOS
ARG TARGETARCH
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build golint -set_exit_status ./...
RUN --mount=type=cache,target=/root/.cache/go-build go vet -v ./...
#RUN mkdir -p dist
#RUN --mount=type=cache,target=/root/.cache/go-build GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -v -x -o dist -ldflags "-X github.com/nrocco/ide/cmd.version=${BUILD_VERSION} -X github.com/nrocco/ide/cmd.commit=${BUILD_COMMIT} -X github.com/nrocco/ide/cmd.date=${BUILD_DATE}"
RUN --mount=type=cache,target=/root/.cache/go-build go test -v -short ./...

# FROM scratch
# COPY --from=builder /app/dist/ /
