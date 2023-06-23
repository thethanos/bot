package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (h *Handler) GetCities(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	cities, err := h.dbAdapter.GetCities("")
	if err != nil {
		h.logger.Error("server::GetCities::GetCities", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		return
	}

	cityList, err := json.Marshal(&cities)
	if err != nil {
		h.logger.Error("server::GetCities::Marshal", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Write(cityList)
	h.logger.Info("Response sent")
}

func (h *Handler) GetCategories(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	categories, err := h.dbAdapter.GetCategories("")
	if err != nil {
		h.logger.Error("server::GetCategories::GetCategories", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		return
	}

	categoryList, err := json.Marshal(&categories)
	if err != nil {
		h.logger.Error("server::GetCategories::Marshal")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Write(categoryList)
	h.logger.Info("Response sent")
}

func (h *Handler) GetServices(rw http.ResponseWriter, req *http.Request) {
	h.logger.Info("Request received: %s", req.URL)

	query := req.URL.Query()
	categoryId := query.Get("category_id")

	services, err := h.dbAdapter.GetServices(categoryId, "")
	if err != nil {
		h.logger.Error("server::GetServices::GetServices")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		return
	}

	serviceList, err := json.Marshal(&services)
	if err != nil {
		h.logger.Error("server::GetServices::Marshal", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Write(serviceList)
	h.logger.Info("Response sent")
}

func (h *Handler) GetMastersList(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	query := req.URL.Query()
	cityId := query.Get("city_id")
	serviceId := query.Get("service_id")

	masters, err := h.dbAdapter.GetMasters(cityId, serviceId)
	if err != nil {
		h.logger.Error("server::GetMastersList::GetMasters", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		return
	}

	template, err := webapp.GenerateWebPage("Выбор мастера", masters)
	if err != nil {
		h.logger.Error("server::GetMastersList::GenerateWebPage", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	rw.Write(template)
	h.logger.Info("Response sent")
}

/*
	func (h *Handler) GetMasterPreview(rw http.ResponseWriter, req *http.Request) {
		h.logger.Infof("Request received: %s", req.URL)

		query := req.URL.Query()
		master_id := query.Get("master")

		master, err := h.dbAdapter.GetMasterPreview(master_id)
		if err != nil {
			h.logger.Error("server::GetMasterPreview::GetMasterPreview", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			return
		}

		template, err := webapp.GenerateWebPage("Предпросмотр", []*entities.Master{master})
		if err != nil {
			h.logger.Error("server::GetMastersList::GenerateWebPage", err)
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			return
		}

		rw.WriteHeader(http.StatusOK)
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.Write(template)
		h.logger.Info("Response sent")
	}
*/
func (h *Handler) SaveMasterRegForm(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::SaveMasterRegForm::ReadAll", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		return
	}

	regForm := &entities.MasterRegForm{}
	if err := json.Unmarshal(body, regForm); err != nil {
		h.logger.Error("server::SaveMasterRegForm::Unmarshal", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Access-Control-Allow-Origin", "*")
	}

	id, err := h.dbAdapter.SaveMasterRegForm(regForm)
	if err != nil {
		h.logger.Error("server::SaveMasterRegForm::SaveMasterRegForm", err)
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s"}`, id)))
	h.logger.Info("Response sent")
}

func (h *Handler) CommitMaster(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Access-Control-Allow-Origin", "*")
	h.logger.Info("Response sent")
}
