package main

import (
	"context"
	"multimessenger_bot/internal/bot"
	"multimessenger_bot/internal/config"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/logger"
	ma "multimessenger_bot/internal/messenger_adapter"
	srv "multimessenger_bot/internal/server"
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

	logger := logger.NewLogger()

	cfg, err := config.Load("config.toml")
	if err != nil {
		logger.Error("main::config::Load", err)
		return
	}

	dbAdapter, _, err := db_adapter.NewDbAdapter(logger)
	if err != nil {
		logger.Error("main::db_adapter::NewDbAdapter", err)
		return
	}

	if err := dbAdapter.AutoMigrate(); err != nil {
		logger.Error("main::db_adapter::AutoMigrate", err)
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

	server, err := srv.NewServer(logger, dbAdapter)
	if err != nil {
		logger.Error("main::server::NewServer", err)
		return
	}

	go func() {
		if err := server.ListenAndServeTLS("dev-full.crt", "dev-key.key"); err != nil {
			logger.Fatal("main::server::ListenAndServeTLS", err)
		}
	}()

	signalHandler := setupSignalHandler()
	<-signalHandler

	bot.Shutdown()
	server.Shutdown(context.Background())
}

func setupSignalHandler() chan os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	return ch
}
