FROM golang:1.18-alpine AS base
MAINTAINER Arthur Aslanyan <arthur.e.aslanyan@gmail.com>

ENV CGO_ENABLED=0
WORKDIR /src
COPY go.* ./
RUN go mod download

FROM base AS test
COPY . .
RUN go test -coverprofile=coverage.txt -covermode=atomic -v ./...

FROM golangci/golangci-lint:v1.46-alpine AS lint-base

FROM base as lint
RUN --mount=target=. \
    --mount=from=lint-base,src=/usr/bin/golangci-lint,target=/usr/bin/golangci-lint \
    golangci-lint run

FROM scratch AS coverage-test
COPY --from=test /src/coverage.txt /

FROM base AS build
RUN --mount=target=. \
    go build -o /out/webhook -v ./cmd/webhook

FROM scratch AS bin
COPY --from=build /out/webhook /
