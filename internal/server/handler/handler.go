package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/entities"
	"multimessenger_bot/internal/logger"
	"multimessenger_bot/internal/webapp"
	"net/http"
	"os"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	logger    logger.Logger
	dbAdapter *db_adapter.DbAdapter
}

func NewHandler(logger logger.Logger, dbAdapter *db_adapter.DbAdapter) *Handler {
	return &Handler{
		logger:    logger,
		dbAdapter: dbAdapter,
	}
}

// @Summary Get cities
// @Description Get all available cities
// @Tags City
// @Accept json
// @Produce json
// @Success 200 {array} entities.City
// @Failure 500 {string} string "Error message"
// @Router /cities [get]
func (h *Handler) GetCities(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	cities, err := h.dbAdapter.GetCities("")
	if err != nil {
		h.logger.Error("server::GetCities::GetCities", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	cityList, err := json.Marshal(&cities)
	if err != nil {
		h.logger.Error("server::GetCities::Marshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(cityList)
	h.logger.Info("Response sent")
}

// @Summary Get categories
// @Description Get all available service categories
// @Tags Service
// @Acept json
// @Produce json
// @Success 200 {array} entities.ServiceCategory
// @Failure 500 {string} string "Error message"
// @Router /categories [get]
func (h *Handler) GetCategories(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	categories, err := h.dbAdapter.GetCategories("")
	if err != nil {
		h.logger.Error("server::GetCategories::GetCategories", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	categoryList, err := json.Marshal(&categories)
	if err != nil {
		h.logger.Error("server::GetCategories::Marshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(categoryList)
	h.logger.Info("Response sent")
}

// @Summary Get services
// @Description Get all available services, filters by category_id if provided
// @Tags Service
// @Param category_id query string false "ID of the service category"
// @Accept json
// @Produce json
// @Success 200 {array} entities.Service
// @Failure 500 {string} string "Error message"
// @Router /services [get]
func (h *Handler) GetServices(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	query := req.URL.Query()
	categoryId := query.Get("category_id")

	services, err := h.dbAdapter.GetServices(categoryId, "")
	if err != nil {
		h.logger.Error("server::GetServices::GetServices", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	serviceList, err := json.Marshal(&services)
	if err != nil {
		h.logger.Error("server::GetServices::Marshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	rw.Write(serviceList)
	h.logger.Info("Response sent")
}

// @Summary Get masters
// @Description Get all available masters for the selected city and the service
// @Tags Master
// @Param page query string true "Page number for pagination"
// @Param limit query string true "Limit number for pagination"
// @Param city_id query string true "ID of the selected city"
// @Param service_id query string true "ID of the seleted service"
// @Accept json
// @Produce json
// @Success 200 {array} entities.Master
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters [get]
func (h *Handler) GetMasters(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	query := req.URL.Query()
	pageParam := query.Get("page")
	page, err := strconv.Atoi(pageParam)
	if err != nil {
		h.logger.Error("server::GetMasters::Atoi", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	limitParam := query.Get("limit")
	limit, err := strconv.Atoi(limitParam)
	if err != nil {
		h.logger.Error("server::GetMasters::Atoi", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	cityId := query.Get("city_id")
	serviceId := query.Get("service_id")

	masters, err := h.dbAdapter.GetMasters(cityId, serviceId, page, limit)
	if err != nil {
		h.logger.Error("server::GetMasters::GetMasters", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	mastersTemplates := make([]string, 0)
	for _, master := range masters {
		template, err := webapp.GenerateMasterCard(master)
		if err != nil {
			h.logger.Error("server::GetMasters::GenerateMassterCard", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		mastersTemplates = append(mastersTemplates, template)
	}

	mastersResp, err := json.Marshal(mastersTemplates)
	if err != nil {
		h.logger.Error("server::GetMasters::Marshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(mastersResp)
	h.logger.Info("Response sent")
}

// @Summary Save master registration form
// @Description Save registration form of a new master
// @Tags Master
// @Param form body entities.MasterRegForm true "Registration form"
// @Accept json
// @Produce json
// @Success 201 {string} string "Success message"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /msters [post]
func (h *Handler) SaveMasterRegForm(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::SaveMasterRegForm::ReadAll", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	regForm := &entities.MasterRegForm{}
	if err := json.Unmarshal(body, regForm); err != nil {
		h.logger.Error("server::SaveMasterRegForm::Unmarshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	validator := validator.New()
	if err := validator.Struct(regForm); err != nil {
		h.logger.Error("server::SaveMasterRegForm::Struct", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.dbAdapter.SaveMasterRegForm(regForm)
	if err != nil {
		h.logger.Error("server::SaveMasterRegForm::SaveMasterRegForm", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s" }`, id)))
	h.logger.Info("Response sent")
}

func (h *Handler) SaveMasterImage(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	params := mux.Vars(req)
	masterID := params["master_id"]

	req.ParseMultipartForm(10 << 20)
	formFile, meta, err := req.FormFile("image")
	if err != nil {
		h.logger.Error("server::SaveMasterImage::FormFile", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer formFile.Close()

	if err := os.MkdirAll(fmt.Sprintf("./webapp/pages/images/%s", masterID), os.ModePerm); err != nil {
		h.logger.Error("server::SaveMasterImage::MkdirAll", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	image, err := os.Create(fmt.Sprintf("./webapp/pages/images/%s/%s", masterID, meta.Filename))
	if err != nil {
		h.logger.Error("server::SaveMasterImage::Create", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	imageBytes, err := ioutil.ReadAll(formFile)
	if err != nil {
		h.logger.Error("server::SaveMasterImage::ReadAll", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := image.Write(imageBytes); err != nil {
		h.logger.Error("server::SaveMasterImage::Write", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	h.logger.Info("Response sent")
}

func (h *Handler) ApproveMaster(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	params := mux.Vars(req)
	masterID := params["master_id"]

	masterForm, err := h.dbAdapter.GetMasterRegForm(masterID)
	if err != nil {
		h.logger.Error("server::ApproveMaster::GetMasterRegForm", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := h.dbAdapter.SaveMaster(masterForm); err != nil {
		h.logger.Error("server::ApproveMaster::SaveMaster", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	h.logger.Info("Response sent")
}
