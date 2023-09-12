package main

import (
	"fmt"
	"multimessenger_bot/internal/bot"
	"multimessenger_bot/internal/config"
	"multimessenger_bot/internal/dbadapter"
	"multimessenger_bot/internal/logger"
	ma "multimessenger_bot/internal/msgadapter"
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
