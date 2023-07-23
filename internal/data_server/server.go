package data_server

import (
	"errors"
	"fmt"
	"multimessenger_bot/internal/config"
	handler "multimessenger_bot/internal/data_server/handler"
	corsMiddleware "multimessenger_bot/internal/data_server/middleware"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/logger"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func NewDataServer(logger logger.Logger, cfg *config.Config, dbAdapter *db_adapter.DbAdapter) (*http.Server, error) {

	handler := handler.NewHandler(logger, cfg, dbAdapter)
	docHandler := middleware.Redoc(middleware.RedocOpts{SpecURL: "swagger.yaml"}, nil)

	router := mux.NewRouter()
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/cities", handler.GetCities)
	getRouter.HandleFunc("/services/categories", handler.GetServiceCategories)
	getRouter.HandleFunc("/services", handler.GetServices)
	getRouter.HandleFunc("/masters", handler.GetMasters)
	getRouter.PathPrefix("/webapp").Handler(http.FileServer(http.Dir("/multimessenger_bot")))
	getRouter.Handle("/docs", docHandler)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("/multimessenger_bot/docs")))

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/cities", handler.SaveCity)
	postRouter.HandleFunc("/services/categories", handler.SaveServiceCategory)
	postRouter.HandleFunc("/services", handler.SaveService)
	postRouter.HandleFunc("/masters", handler.SaveMasterRegForm)
	postRouter.HandleFunc("/masters/images/{master_id}", handler.SaveMasterImage)
	postRouter.HandleFunc("/masters/approve/{master_id}", handler.ApproveMaster)

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
		Handler: corsMiddleware.CorsMiddlware(router),
		Addr:    addr,
	}, nil
}
