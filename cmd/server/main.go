package main

import (
	"bot/internal/config"
	"bot/internal/dbadapter"
	"bot/internal/logger"
	srv "bot/internal/server"
	"context"
	"fmt"
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

	DBAdapter, err := dbadapter.NewDbAdapter(logger, cfg)
	if err != nil {
		logger.Error("main::dbadapter::NewDBAdapter", err)
		return
	}

	if err := DBAdapter.AutoMigrate(); err != nil {
		logger.Error("main::dbadapter::AutoMigrate", err)
		return
	}

	server, err := srv.NewServer(logger, cfg, DBAdapter)
	if err != nil {
		logger.Error("main::server::NewServer", err)
		return
	}

	go func() {
		if err := server.ListenAndServeTLS("dev-full.crt", "dev-key.key"); err != nil {
			logger.Fatal("main::server::ListenAndServe", err)
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
