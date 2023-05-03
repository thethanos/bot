package main

import (
	"fmt"
	"multimessenger_bot/internal/config"
	"multimessenger_bot/internal/google_calendar"
	"multimessenger_bot/internal/telegram"
	"multimessenger_bot/internal/whatsapp"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {

	gCalendarClient, err := google_calendar.NewGoogleCalendarClient()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(gCalendarClient)

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

	tgClient, _ := telegram.NewTelegramClient(cfg)
	tgClient.Connect()

	waClient, _ := whatsapp.NewWhatsAppClient(nil, cfg, container)
	//waClient.Connect()

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
