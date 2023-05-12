package server

import (
	"encoding/json"
	"io/ioutil"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/entities"
	"net/http"

	"go.uber.org/zap"
)

type Handler struct {
	logger    *zap.SugaredLogger
	dbAdapter *db_adapter.DbAdapter
}

func NewHandler(logger *zap.SugaredLogger, dbAdapter *db_adapter.DbAdapter) *Handler {
	return &Handler{
		logger:    logger,
		dbAdapter: dbAdapter,
	}
}

func (h *Handler) SaveNewCity(rw http.ResponseWriter, req *http.Request) {

}

func (h *Handler) SaveNewService(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received %s", req.URL)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::Handler::SaveNewService", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	service := &entities.Service{}
	if err := json.Unmarshal(body, service); err != nil {
		h.logger.Error("server::Handler::SaveNewService", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.dbAdapter.SaveNewService(service.Name); err != nil {
		h.logger.Error("server::Handler::SaveNewService", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
