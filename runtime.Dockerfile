FROM golang:1.20.3-buster

WORKDIR /multimessenger_bot
COPY . ./
RUN go get -u gorm.io/gorm && go get -u gorm.io/driver/sqlite

RUN go build -o multimessenger_bot main.go

CMD ["./multimessenger_bot"]