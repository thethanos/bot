package bot

import (
	"bot/internal/config"
	dbmock "bot/internal/dbadapter/mock"
	"bot/internal/logger"
	client "bot/internal/messenger_client"
	clmock "bot/internal/messenger_client/mock"
	ma "bot/internal/msgadapter"
	"context"
	"fmt"
	"sync"
	"testing"
)

var messages []*ma.Message

const message_count = 10000

func init() {
	messages = make([]*ma.Message, 0)
	for i := 0; i < message_count; i++ {
		messages = append(messages, &ma.Message{
			Text:   "/start",
			Source: ma.TELEGRAM,
			Type:   ma.TEXT,
			UserID: fmt.Sprintf("%d", i),
			Data:   &ma.MessageData{},
		})
	}

	for i := 0; i < message_count; i++ {
		messages = append(messages, &ma.Message{
			Text:   "город",
			Source: ma.TELEGRAM,
			Type:   ma.TEXT,
			UserID: fmt.Sprintf("%d", i),
			Data:   &ma.MessageData{},
		})
	}

	for i := 0; i < message_count; i++ {
		messages = append(messages, &ma.Message{
			Text:   "Tel-Aviv",
			Source: ma.TELEGRAM,
			Type:   ma.TEXT,
			UserID: fmt.Sprintf("%d", i),
			Data:   &ma.MessageData{},
		})
	}

	for i := 0; i < message_count; i++ {
		messages = append(messages, &ma.Message{
			Text:   "Face",
			Source: ma.TELEGRAM,
			Type:   ma.TEXT,
			UserID: fmt.Sprintf("%d", i),
			Data:   &ma.MessageData{},
		})
	}
}

func BenchmarkMessageHandling(b *testing.B) {

	cfg := &config.Config{}
	logger := logger.NewLogger()
	DBAdapter := &dbmock.DBAdapterMock{}

	recvMsgChan := make(chan *ma.Message, 10)
	resultChan := make(chan *ma.Message)

	wgTest := &sync.WaitGroup{}
	wgTest.Add(message_count * 4)
	clientMock := &clmock.ClientMock{Result: resultChan, Wg: wgTest}

	bot, _ := NewBot(logger, cfg, []client.ClientInterface{clientMock}, DBAdapter, recvMsgChan)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for _, msg := range messages {
			recvMsgChan <- msg
		}
		wgTest.Wait()
		cancel()
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)
	bot.Run(ctx, &wg)
	wg.Wait()
}
