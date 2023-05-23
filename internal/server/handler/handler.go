package server

import (
	"bytes"
	"fmt"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/entities"
	"net/http"
	"text/template"

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

	_, err := h.dbAdapter.GetMasters(cityId, serviceId)
	if err != nil {
		h.logger.Error("server::Handler::GetMastersList::GetMasters", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	masters := []*entities.Master{
		{
			Name:        "Masha",
			Img:         "masters/images/maria_ernandes/1.png",
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		},
		{
			Name:        "Pasha",
			Img:         "masters/images/maria_ernandes/1.png",
			Description: "Lorem ipsum dolor sit amet, consectetur adipiscing elit.",
		},
	}

	allFiles := []string{"content.tmpl", "footer.tmpl", "header.tmpl", "page.tmpl"}

	var allPaths []string
	for _, tmpl := range allFiles {
		allPaths = append(allPaths, "./webapp/masters/templates/"+tmpl)
	}

	templates := template.Must(template.New("").ParseFiles(allPaths...))

	var processed bytes.Buffer
	if err := templates.ExecuteTemplate(&processed, "page", masters); err != nil {
		h.logger.Error("server::Handler::GetMastersList::ExecuteTemplate", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(rw, string(processed.Bytes()))
}
