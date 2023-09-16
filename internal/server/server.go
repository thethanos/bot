package server

import (
	"bot/internal/config"
	"bot/internal/dbadapter"
	"bot/internal/logger"
	handler "bot/internal/server/handler"
	corsMiddleware "bot/internal/server/middleware"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	"github.com/gorilla/mux"
)

func NewServer(logger logger.Logger, cfg *config.Config, DBAdapter *dbadapter.DBAdapter) (*http.Server, error) {

	handler := handler.NewHandler(logger, cfg, DBAdapter)
	docHandler := middleware.Redoc(middleware.RedocOpts{SpecURL: "swagger.yaml"}, nil)

	router := mux.NewRouter()
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/cities", handler.GetCities)
	getRouter.HandleFunc("/services/categories", handler.GetServiceCategories)
	getRouter.HandleFunc("/services", handler.GetServices)
	getRouter.HandleFunc("/masters", handler.GetMasters)
	getRouter.PathPrefix("/images").Handler(http.FileServer(http.Dir("/bot")))
	getRouter.Handle("/docs", docHandler)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("/bot/docs")))

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/cities", handler.SaveCity)
	postRouter.HandleFunc("/services/categories", handler.SaveServiceCategory)
	postRouter.HandleFunc("/services", handler.SaveService)
	postRouter.HandleFunc("/masters", handler.SaveMasterRegForm)
	postRouter.HandleFunc("/masters/images/{master_id}", handler.SaveMasterImage)
	postRouter.HandleFunc("/masters/approve/{master_id}", handler.ApproveMaster)

	deleteHandler := router.Methods(http.MethodDelete).Subrouter()
	deleteHandler.HandleFunc("/cities/{city_id}", handler.DeleteCity)
	deleteHandler.HandleFunc("/services/categories/{category_id}", handler.DeleteServCategory)
	deleteHandler.HandleFunc("/services/{service_id}", handler.DeleteService)
	deleteHandler.HandleFunc("/masters/{master_id}", handler.DeleteMaster)

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
