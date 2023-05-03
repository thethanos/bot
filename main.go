package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"whatsapp_bot/internal/config"
	"whatsapp_bot/internal/telegram_client"
	"whatsapp_bot/internal/whatsapp_client"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	cfg, err := config.Load("config.toml")
	if err != nil {
		fmt.Print(err)
		return
	}

	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	//clientLog := waLog.Stdout("Client", "DEBUG", true)

	tgClient, _ := telegram_client.NewTelegramClient(cfg)
	tgClient.Connect()

	waClient, _ := whatsapp_client.NewWhatsAppClient(nil, cfg, container)
	waClient.Connect()

	signalHandler := setupSignalHandler()
	<-signalHandler

	tgClient.Disconnect()
	waClient.Disconnect()
}

func setupSignalHandler() chan os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	return ch
}
