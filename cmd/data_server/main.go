package main

import (
	"context"
	"fmt"
	"multimessenger_bot/internal/config"
	srv "multimessenger_bot/internal/data_server"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {

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

	if err := dbAdapter.AutoMigrate(); err != nil {
		logger.Error("main::db_adapter::AutoMigrate", err)
		return
	}

	server, err := srv.NewDataServer(logger, cfg, dbAdapter)
	if err != nil {
		logger.Error("main::data_server::NewDataServer", err)
		return
	}

	go func() {
		if err := server.ListenAndServeTLS("dev-full.crt", "dev-key.key"); err != nil {
			logger.Fatal("main::data_server::ListenAndServe", err)
		}
	}()

	signalHandler := setupSignalHandler()
	<-signalHandler

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Error("main::data_server::Shutdown", err)
	}
}

func setupSignalHandler() chan os.Signal {
	size := 2
	ch := make(chan os.Signal, size)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	return ch
}
