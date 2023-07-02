FROM golang:1.20.3-buster

WORKDIR /multimessenger_bot
COPY . ./

RUN go install -v github.com/magefile/mage@latest
RUN mage build

CMD ["./multimessenger_bot"]