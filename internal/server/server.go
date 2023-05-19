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
	getRouter.HandleFunc("/master", handler.GetMastersList)
	getRouter.PathPrefix("/").Handler(http.FileServer(http.Dir("./page")))

	server := &http.Server{
		Handler: router,
		Addr:    ":443",
	}

	return server, nil
}
