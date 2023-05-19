package server

import (
	"bytes"
	"fmt"
	"multimessenger_bot/internal/db_adapter"
	"net/http"
	"text/template"

	"go.uber.org/zap"
)

type product struct {
	Img         string
	Name        string
	Stars       float64
	Description string
}

func subtr(a, b float64) float64 {
	return a - b
}

func list(e ...float64) []float64 {
	return e
}

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

	for idx, master := range masters {
		master.Stars = float64(idx)
		master.Description = "Lorem ipsum dolor sit amet, consectetur adipiscing elit."
		master.Img = "images/1.png"
	}

	allFiles := []string{"content.tmpl", "footer.tmpl", "header.tmpl", "page.tmpl"}

	var allPaths []string
	for _, tmpl := range allFiles {
		allPaths = append(allPaths, "./page/templates/"+tmpl)
	}

	templates := template.Must(template.New("").Funcs(template.FuncMap{"subtr": subtr, "list": list}).ParseFiles(allPaths...))

	var processed bytes.Buffer
	templates.ExecuteTemplate(&processed, "page", masters)

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(rw, string(processed.Bytes()))
}
