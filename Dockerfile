FROM golang:1.18-alpine AS base
MAINTAINER Arthur Aslanyan <arthur.e.aslanyan@gmail.com>

ENV CGO_ENABLED=0
WORKDIR /src
COPY go.* .
RUN go mod download

FROM base AS test
COPY . .
RUN go test -v ./...

FROM golangci/golangci-lint:v1.46-alpine AS lint-base

FROM base as lint
COPY --from=lint-base /usr/bin/golangci-lint /usr/bin/golangci-lint
RUN --mount=target=. \
    golangci-lint run
