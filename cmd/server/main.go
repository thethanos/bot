package main

import (
	"context"
	"fmt"
	"multimessenger_bot/internal/config"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/logger"
	srv "multimessenger_bot/internal/server"
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

	server, err := srv.NewServer(logger, cfg, dbAdapter)
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

	if err := server.Shutdown(context.Background()); err != nil {
		logger.Error("main::server::Shutdown", err)
	}
}

func setupSignalHandler() chan os.Signal {
	size := 2
	ch := make(chan os.Signal, size)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	return ch
}
