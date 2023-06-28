package server

import (
	"errors"
	"fmt"
	"multimessenger_bot/internal/config"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/logger"
	handler "multimessenger_bot/internal/server/handler"
	"net/http"

	"github.com/gorilla/mux"
)

func NewServer(logger logger.Logger, cfg *config.Config, dbAdapter *db_adapter.DbAdapter) (*http.Server, error) {

	handler := handler.NewHandler(logger, dbAdapter)

	router := mux.NewRouter()
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/cities", handler.GetCities)
	getRouter.HandleFunc("/categories", handler.GetCategories)
	getRouter.HandleFunc("/services", handler.GetServices)
	getRouter.HandleFunc("/masters", handler.GetMasters)
	getRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("./webapp")))

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/masters", handler.SaveMasterRegForm)
	postRouter.HandleFunc("/masters/images/{master_id}", handler.SaveMasterImage)

	var addr string
	switch cfg.Mode {
	case config.DEBUG:
		addr = fmt.Sprintf(":%d", cfg.DebugPort)
	case config.RELEASE:
		addr = fmt.Sprintf(":%d", cfg.ReleasePort)
	default:
		return nil, errors.New("Run mode is not specified")
	}

	return &http.Server{
		Handler: router,
		Addr:    addr,
	}, nil
}
