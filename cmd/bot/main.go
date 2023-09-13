package main

import (
	"bot/internal/bot"
	"bot/internal/config"
	"bot/internal/dbadapter"
	"bot/internal/logger"
	ma "bot/internal/msgadapter"
	"bot/internal/telegram"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	cfg, err := config.Load("config.toml")
	if err != nil {
		panic(fmt.Sprintf("main::config::Load::%s", err))
	}

	logger := logger.NewLogger(cfg.Mode)

	DBAdapter, err := dbadapter.NewDbAdapter(logger, cfg)
	if err != nil {
		logger.Error("main::dbadapter::NewDbAdapter", err)
		return
	}

	recvMsgChan := make(chan *ma.Message)
	tgClient, _ := telegram.NewTelegramClient(logger, cfg, recvMsgChan)
	//waClient, _ := whatsapp.NewWhatsAppClient(logger, cfg, waContainer, recvMsgChan)

	bot, err := bot.NewBot(logger, []ma.ClientInterface{tgClient}, DBAdapter, recvMsgChan)
	if err != nil {
		logger.Error("main::bot::NewBot", err)
	}
	bot.Run()

	signalHandler := setupSignalHandler()
	<-signalHandler

	bot.Shutdown()
}

func setupSignalHandler() chan os.Signal {
	size := 2
	ch := make(chan os.Signal, size)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	return ch
}
