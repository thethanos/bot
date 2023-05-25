package server

import (
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/entities"
	"multimessenger_bot/internal/webapp"
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

func (h *Handler) GetMastersList(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	query := req.URL.Query()
	cityId := query.Get("city")
	serviceId := query.Get("service")

	masters, err := h.dbAdapter.GetMasters(cityId, serviceId)
	if err != nil {
		h.logger.Error("server::Handler::GetMastersList::GetMasters", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	template, err := webapp.GenerateWebPage("Выбор мастера", masters)
	if err != nil {
		h.logger.Error("server::Handler::GetMastersList::ExecuteTemplate", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.Write(template)
}

func (h *Handler) GetMasterPreview(rw http.ResponseWriter, req *http.Request) {
	template, err := webapp.GenerateWebPage("Предпросмотр", []*entities.Master{{Name: "Test", Image: "https://bot-dev-domain.com/masters/images/maria_ernandes/1.png", Description: "lorem ipsum"}})
	if err != nil {
		h.logger.Error("server::Handler::GetMastersList::ExecuteTemplate", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.Write(template)
}
