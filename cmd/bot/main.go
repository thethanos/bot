package main

import (
	"fmt"
	"multimessenger_bot/internal/bot"
	"multimessenger_bot/internal/config"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/logger"
	ma "multimessenger_bot/internal/messenger_adapter"
	"multimessenger_bot/internal/telegram"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	//gCalendarClient, err := google_calendar.NewGoogleCalendarClient("service_credentials.json")
	//if err != nil {
	//	  fmt.Println(err)
	//	  return
	//}

	cfg, err := config.Load("config.toml")
	if err != nil {
		panic(fmt.Sprintf("main::config::Load::%s", err))
	}

	logger := logger.NewLogger(cfg.Mode)

	dbAdapter, err := db_adapter.NewDbAdapter(logger, cfg)
	if err != nil {
		logger.Error("main::db_adapter::NewDbAdapter", err)
		return
	}

	recvMsgChan := make(chan *ma.Message)
	tgClient, _ := telegram.NewTelegramClient(logger, cfg, recvMsgChan)
	//waClient, _ := whatsapp.NewWhatsAppClient(logger, cfg, waContainer, recvMsgChan)

	bot, err := bot.NewBot(logger, []ma.ClientInterface{tgClient}, dbAdapter, recvMsgChan)
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
