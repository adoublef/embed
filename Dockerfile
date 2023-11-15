# syntax=docker/dockerfile:1

ARG GO_VERSION=1.21
ARG ALPINE_VERSION=3.18

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build
WORKDIR /usr/src

ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# required for go-sqlite3
RUN apk add --no-cache gcc musl-dev

COPY go.* .
RUN go mod download

COPY . .
RUN go build \
    -ldflags "-s -w -extldflags '-static'" \
    -buildvcs=false \
    -o /usr/bin/ ./...

FROM alpine:${ALPINE_VERSION} AS runtime
WORKDIR /opt

ARG EXE_NAME
COPY --from=build /usr/local/bin/${EXE_NAME} ./a

# run migration scripts here?

CMD ["./a", "s"]