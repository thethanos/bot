FROM golang:1.20.3-buster

WORKDIR /multimessenger_bot
COPY . ./

RUN go install -v github.com/magefile/mage@latest && go install -v github.com/swaggo/swag/cmd/swag@latest && go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@v1.53.3
RUN mage -d cmd buildDataServer

CMD ["./server"]