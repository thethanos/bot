FROM golang:1.20.3-buster

RUN go install -v golang.org/x/tools/gopls@latest && \
    go install -v github.com/go-delve/delve/cmd/dlv@latest && \
    go install -v github.com/magefile/mage@latest && \
    go install -v github.com/swaggo/swag/cmd/swag@latest && \
    go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3

WORKDIR /bot
COPY go.mod go.sum ./


RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && apt-get update && apt-get upgrade -y && \
    apt-get install -y apt-transport-https ca-certificates gnupg python3 psmisc nodejs

RUN echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
RUN curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key --keyring /usr/share/keyrings/cloud.google.gpg add -

RUN apt-get update && apt-get install -y google-cloud-cli