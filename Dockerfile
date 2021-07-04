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
RUN go build -tags netgo -o main ./


FROM alpine:latest as app
COPY --from=builder /app/main /usr/bin/main
ENTRYPOINT [ "/usr/bin/main" ]
CMD []