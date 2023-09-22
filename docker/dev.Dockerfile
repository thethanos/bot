FROM golang:1.20.3-buster

RUN go install -v golang.org/x/tools/gopls@latest && \
    go install -v github.com/go-delve/delve/cmd/dlv@latest && \
    go install -v github.com/magefile/mage@latest && \
    go install -v github.com/swaggo/swag/cmd/swag@latest && \
    go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

WORKDIR /bot
COPY go.mod go.sum ./


RUN apt-get update && apt-get upgrade -y
