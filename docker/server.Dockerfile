FROM golang:1.20.3-buster

WORKDIR /multimessenger_bot
COPY . ./

RUN go install -v github.com/magefile/mage@latest && go install -v github.com/swaggo/swag/cmd/swag@latest
RUN mage -d cmd buildServer

CMD ["./server"]