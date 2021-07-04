FROM golangci/golangci-lint:latest as linter
WORKDIR /app
COPY . .
RUN golangci-lint run


FROM golang:1.15 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go version
RUN go mod download
COPY . .
RUN go test -v ./...
RUN go build -tags netgo -o dead-link-finder ./


FROM alpine:latest as app
RUN apk --no-cache add ca-certificates curl openssh-client git
COPY --from=builder /app/dead-link-finder /usr/bin/dead-link-finder