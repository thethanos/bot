FROM golang:1.20.3-buster

RUN go install -v golang.org/x/tools/gopls@latest