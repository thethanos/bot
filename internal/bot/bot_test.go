package bot

import (
	"bot/internal/config"
	dbmock "bot/internal/dbadapter/mock"
	"bot/internal/logger"
	client "bot/internal/messenger_client"
	clmock "bot/internal/messenger_client/mock"
	ma "bot/internal/msgadapter"
	"context"
	"sync"
	"testing"
)

func BenchmarkMessageHandling(b *testing.B) {

	cfg := &config.Config{}
	logger := logger.NewLogger()
	DBAdapter := &dbmock.DBAdapterMock{}

	recvMsgChan := make(chan *ma.Message, 10)
	clientMock := &clmock.ClientMock{}

	bot, _ := NewBot(logger, cfg, []client.ClientInterface{clientMock}, DBAdapter, recvMsgChan)

	wg := sync.WaitGroup{}
	wg.Add(3)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {

		cancel()
	}()

	bot.Run(ctx, &wg)
}
