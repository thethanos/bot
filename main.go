package main

import (
	"os"
	"os/signal"
	"syscall"
	"whatsapp_bot/internal/whatsapp_client"

	_ "github.com/mattn/go-sqlite3"
	"go.mau.fi/whatsmeow/store/sqlstore"
	waLog "go.mau.fi/whatsmeow/util/log"
)

func main() {

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	// Make sure you add appropriate DB connector imports, e.g. github.com/mattn/go-sqlite3 for SQLite
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil {
		panic(err)
	}

	//clientLog := waLog.Stdout("Client", "DEBUG", true)

	waClient := whatsapp_client.NewWhatsAppClient(nil, container)
	waClient.Connect()

	signalHandler := setupSignalHandler()
	<-signalHandler

	waClient.Disconnect()
}

func setupSignalHandler() chan os.Signal {
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	return ch
}
