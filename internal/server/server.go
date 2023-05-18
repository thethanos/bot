package server

import (
	"multimessenger_bot/internal/db_adapter"
	handler "multimessenger_bot/internal/server/handler"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewServer(logger *zap.SugaredLogger, dbAdapter *db_adapter.DbAdapter) (*http.Server, error) {

	handler := handler.NewHandler(logger, dbAdapter)

	router := mux.NewRouter()
	getRouter := router.Methods(http.MethodGet).Subrouter()
	getRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))
	getRouter.HandleFunc("/masters", handler.GetMastersList)

	postRouter := router.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/city", handler.SaveNewCity)
	postRouter.HandleFunc("/service", handler.SaveNewService)

	server := &http.Server{
		Handler: router,
		Addr:    ":443",
	}

	return server, nil
}
