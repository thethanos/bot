package main

import (
	"context"
	"fmt"
	"multimessenger_bot/internal/bot"
	"multimessenger_bot/internal/config"
	"multimessenger_bot/internal/db_adapter"
	ma "multimessenger_bot/internal/messenger_adapter"
	"multimessenger_bot/internal/telegram"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"

	_ "github.com/mattn/go-sqlite3"
)

func main() {

	//gCalendarClient, err := google_calendar.NewGoogleCalendarClient("service_credentials.json")
	//if err != nil {
	//	  fmt.Println(err)
	//	  return
	//}

	//dbLog := waLog.Stdout("Database", "DEBUG", true)
	cfg, err := config.Load("config.toml")
	if err != nil {
		fmt.Print(err)
		return
	}

	dbAdapter, _, err := db_adapter.NewDbAdapter()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := dbAdapter.AutoMigrate(); err != nil {
		fmt.Println(err)
		return
	}

	//dbAdapter.Test()

	//clientLog := waLog.Stdout("Client", "DEBUG", true)

	recvMsgChan := make(chan *ma.Message)
	tgClient, _ := telegram.NewTelegramClient(cfg, recvMsgChan)
	//waClient, _ := whatsapp.NewWhatsAppClient(nil, cfg, waContainer, recvMsgChan)

	bot, _ := bot.NewBot([]ma.ClientInterface{tgClient}, dbAdapter, recvMsgChan)
	bot.Run()

	// Setup new HTTP server mux to handle different paths.
	mux := http.NewServeMux()
	// This serves the home page.
	mux.HandleFunc("/", index("/workspaces/multimessenger_bot/index.html"))
	// This serves our "validation" API, which checks if the input data is valid.
	server := http.Server{
		Handler: mux,
		Addr:    "0.0.0.0:8082",
	}

	go func() {
		if err := server.ListenAndServeTLS("localhost.pem", "localhost-key.pem"); err != nil {
			panic("failed to listen and serve: " + err.Error())
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

var indexTmpl = template.Must(template.ParseFiles("index.html"))

func index(webappURL string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := indexTmpl.ExecuteTemplate(writer, "index.html", struct {
			WebAppURL string
		}{
			WebAppURL: webappURL,
		})
		if err != nil {
			writer.WriteHeader(http.StatusInternalServerError)
			writer.Write([]byte(err.Error()))
		}
	}
}
