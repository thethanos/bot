package main

import (
	"bot/internal/bot"
	"bot/internal/config"
	"bot/internal/dbadapter"
	"bot/internal/logger"
	client "bot/internal/messenger_client"
	"bot/internal/messenger_client/telegram"
	ma "bot/internal/msgadapter"
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	cfg, err := config.Load("config.toml")
	if err != nil {
		panic(fmt.Sprintf("main::config::Load::%s", err))
	}

	logger := logger.NewLogger()

	DBAdapter, err := dbadapter.NewDbAdapter(logger, cfg)
	if err != nil {
		logger.Error("main::dbadapter::NewDbAdapter", err)
		return
	}

	recvMsgChan := make(chan *ma.Message, cfg.RcvBufSize)
	tgClient, _ := telegram.NewTelegramClient(logger, cfg, recvMsgChan)
	//waClient, _ := whatsapp.NewWhatsAppClient(logger, cfg, waContainer, recvMsgChan)

	bot, err := bot.NewBot(logger, cfg, []client.ClientInterface{tgClient}, DBAdapter, recvMsgChan)
	if err != nil {
		logger.Error("main::bot::NewBot", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(2)

	bot.Run(ctx, &wg)

	signalHandler := setupSignalHandler()
	<-signalHandler

	cancel()
	wg.Wait()
	bot.Shutdown()
}

func setupSignalHandler() chan os.Signal {
	size := 2
	ch := make(chan os.Signal, size)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	return ch
}
