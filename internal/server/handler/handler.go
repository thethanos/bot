package server

import (
	"encoding/json"
	"fmt"
	"io"
	"multimessenger_bot/internal/config"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/entities"
	"multimessenger_bot/internal/logger"
	"multimessenger_bot/internal/webapp"
	"net/http"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type Handler struct {
	logger    logger.Logger
	cfg       *config.Config
	dbAdapter *db_adapter.DbAdapter
}

func NewHandler(logger logger.Logger, cfg *config.Config, dbAdapter *db_adapter.DbAdapter) *Handler {
	return &Handler{
		logger:    logger,
		cfg:       cfg,
		dbAdapter: dbAdapter,
	}
}

// @Summary Get cities
// @Description Get all available cities
// @Tags City
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Limit of items for pagination"
// @Accept json
// @Produce json
// @Success 200 {array} entities.City
// @Failure 500 {string} string "Error message"
// @Router /cities [get]
func (h *Handler) GetCities(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	query := req.URL.Query()
	page, err := getParamInt(query.Get("page"), 0)
	if err != nil {
		h.logger.Error("server::GetCities::getParamInt", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := getParamInt(query.Get("limit"), -1)
	if err != nil {
		h.logger.Error("server::GetCities::getParamInt", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	cities, err := h.dbAdapter.GetCities("", page, limit)
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
	if _, err := rw.Write(cityList); err != nil {
		h.logger.Error("server::GetCities::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Get service categories
// @Description Get all available service categories
// @Tags Service
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Limit of items for pagination"
// @Acept json
// @Produce json
// @Success 200 {array} entities.ServiceCategory
// @Failure 500 {string} string "Error message"
// @Router /services/categories [get]
func (h *Handler) GetServiceCategories(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	query := req.URL.Query()
	page, err := getParamInt(query.Get("page"), 0)
	if err != nil {
		h.logger.Error("server::GetCities::getParamInt", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := getParamInt(query.Get("limit"), -1)
	if err != nil {
		h.logger.Error("server::GetCities::getParamInt", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	categories, err := h.dbAdapter.GetServiceCategories("", page, limit)
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
	if _, err := rw.Write(categoryList); err != nil {
		h.logger.Error("server::GetServiceCategories::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Get services
// @Description Get all available services, filters by category_id if provided
// @Tags Service
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Limit of items for pagination"
// @Param category_id query string false "ID of the service category"
// @Accept json
// @Produce json
// @Success 200 {array} entities.Service
// @Failure 500 {string} string "Error message"
// @Router /services [get]
func (h *Handler) GetServices(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	query := req.URL.Query()
	page, err := getParamInt(query.Get("page"), 0)
	if err != nil {
		h.logger.Error("server::GetCities::getParamInt", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := getParamInt(query.Get("limit"), -1)
	if err != nil {
		h.logger.Error("server::GetCities::getParamInt", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	categoryId := query.Get("category_id")

	services, err := h.dbAdapter.GetServices(categoryId, "", page, limit)
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
	if _, err := rw.Write(serviceList); err != nil {
		h.logger.Error("server::GetServices::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Get masters
// @Description Get all available masters for the selected city and the service
// @Tags Master
// @Param page query string false "Page number for pagination"
// @Param limit query string false "Limit of items for pagination"
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
	page, err := getParamInt(query.Get("page"), 0)
	if err != nil {
		h.logger.Error("server::GetMasters::getParamInt", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := getParamInt(query.Get("limit"), -1)
	if err != nil {
		h.logger.Error("server::GetMasters::getParamInt", err)
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

	mastersResp, err := json.Marshal(masters)
	if err != nil {
		h.logger.Error("server::GetMasters::Marshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(mastersResp); err != nil {
		h.logger.Error("server::GetMasters::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Get masters in html
// @Description Get all available masters wrapped up in the html blocks
// @Tags Master
// @Param page query string true "Page number for pagination"
// @Param limit query string true "Limit of items for pagination"
// @Param city_id query string true "ID of the selected city"
// @Param service_id query string true "ID of the seleted service"
// @Accept json
// @Produce json
// @Success 200 {string} string "Collection of HTML blocks"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters/html [get]
func (h *Handler) GetMastersHTML(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	query := req.URL.Query()
	page, err := getParamInt(query.Get("page"), 0)
	if err != nil {
		h.logger.Error("server::GetMastersHTML::getParamInt", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	limit, err := getParamInt(query.Get("limit"), -1)
	if err != nil {
		h.logger.Error("server::GetMastersHTML::getParamInt", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	cityId := query.Get("city_id")
	serviceId := query.Get("service_id")

	masters, err := h.dbAdapter.GetMasters(cityId, serviceId, page, limit)
	if err != nil {
		h.logger.Error("server::GetMastersHTML::GetMasters", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	mastersTemplates := make([]string, 0)
	for _, master := range masters {
		template, err := webapp.GenerateMasterCard(master)
		if err != nil {
			h.logger.Error("server::GetMastersHTML::GenerateMassterCard", err)
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		mastersTemplates = append(mastersTemplates, template)
	}

	mastersResp, err := json.Marshal(mastersTemplates)
	if err != nil {
		h.logger.Error("server::GetMastersHTML::Marshal", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
	if _, err := rw.Write(mastersResp); err != nil {
		h.logger.Error("server::GetMastersHTML::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Save city
// @Description Save a new city in the system
// @Tags City
// @Param name body Name true "City name"
// @Accept json
// @Produce json
// @Success 201 {object} ID "ID of the new city"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /cities [post]
func (h *Handler) SaveCity(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::SaveCity::ReadAll")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	city := &entities.City{}
	if err := json.Unmarshal(body, city); err != nil {
		h.logger.Error("server::SaveCity::Unmarshal")
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.dbAdapter.SaveCity(city.Name)
	if err != nil {
		h.logger.Error("server::SaveCity::SaveCity", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s" }`, id))); err != nil {
		h.logger.Error("server::SaveCity::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Save service category
// @Description Save a new service category in the system
// @Tags Service
// @Param name body Name true "Service category name"
// @Accept json
// @Produce json
// @Success 201 {object} ID "ID of the new service category"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /services/categories [post]
func (h *Handler) SaveServiceCategory(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::SaveServiceCategory::ReadAll", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	serviceCategory := &entities.ServiceCategory{}
	if err := json.Unmarshal(body, serviceCategory); err != nil {
		h.logger.Error("server::SaveServiceCategory::Unmarshal", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.dbAdapter.SaveServiceCategory(serviceCategory.Name)
	if err != nil {
		h.logger.Error("server::SaveServiceCategory::SaveServiceCategory", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s" }`, id))); err != nil {
		h.logger.Error("server::SaveServiceCategory::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Save service
// @Description Save a new service in the system
// @Tags Service
// @Param service body entities.Service true "New service"
// @Accept json
// @Produce json
// @Success 201 {object} ID "ID of the new service"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /services [post]
func (h *Handler) SaveService(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::SaveService::ReadAll", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	service := &entities.Service{}
	if err := json.Unmarshal(body, service); err != nil {
		h.logger.Error("server::SaveService::Unmarshal", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	id, err := h.dbAdapter.SaveService(service.Name, service.CategoryID)
	if err != nil {
		h.logger.Error("server::SaveService::SaveService", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s" }`, id))); err != nil {
		h.logger.Error("server::SaveService::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Save master registration form
// @Description Save registration form of a new master
// @Tags Master
// @Param form body entities.MasterRegForm true "Registration form"
// @Accept json
// @Produce json
// @Success 201 {object} ID "ID of the new master"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters [post]
func (h *Handler) SaveMasterRegForm(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	body, err := io.ReadAll(req.Body)
	if err != nil {
		h.logger.Error("server::SaveMasterRegForm::ReadAll", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	regForm := &entities.MasterRegForm{}
	if err := json.Unmarshal(body, regForm); err != nil {
		h.logger.Error("server::SaveMasterRegForm::Unmarshal", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
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
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s" }`, id))); err != nil {
		h.logger.Error("server::SaveMasterRegForm::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Save master's image
// @Description Save the image that was attached to the registration form
// @Tags Master
// @Param master_id path string true "ID of a master, whose picture is uploaded"
// @Param file formData file true "Image to upload"
// @Accept multipart/form-data
// @Produce json
// @Success 201 {object} URL "URL of the saved picture"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters/images/{master_id} [post]
func (h *Handler) SaveMasterImage(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	params := mux.Vars(req)
	masterID := params["master_id"]
	if len(masterID) == 0 {
		h.logger.Error("server::SaveMasterImage::params[]", "no masterID")
		http.Error(rw, "no masterID", http.StatusBadRequest)
		return
	}

	if err := req.ParseMultipartForm(10 << 20); err != nil {
		h.logger.Error("server::SaveMasterImage::ParseMultipartForm", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	formFile, meta, err := req.FormFile("file")
	if err != nil {
		h.logger.Error("server::SaveMasterImage::FormFile", err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
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

	imageBytes, err := io.ReadAll(formFile)
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

	imgUrl := fmt.Sprintf("%s/%s/%s", h.cfg.ImagePrefix, masterID, meta.Filename)

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "url" : "%s" }`, imgUrl))); err != nil {
		h.logger.Error("server::SaveMasterImage::Write", err)
		return
	}
	h.logger.Info("Response sent")
}

// @Summary Approve master
// @Description Approve and save master in the system
// @Tags Master
// @Param master_id path string true "ID of the approved master"
// @Accept json
// @Produce json
// @Success 201 {object} ID "ID of the approved master"
// @Failure 400 {string} string "Error message"
// @Failure 500 {string} string "Error message"
// @Router /masters/approve/{maser_id} [post]
func (h *Handler) ApproveMaster(rw http.ResponseWriter, req *http.Request) {
	h.logger.Infof("Request received: %s", req.URL)

	params := mux.Vars(req)
	masterID := params["master_id"]
	if len(masterID) == 0 {
		h.logger.Error("server::ApproveMaster::params[]", "no masterID")
		http.Error(rw, "no masterID", http.StatusBadRequest)
		return
	}

	masterForm, err := h.dbAdapter.GetMasterRegForm(masterID)
	if err != nil {
		h.logger.Error("server::ApproveMaster::GetMasterRegForm", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, err := h.dbAdapter.SaveMaster(masterForm); err != nil {
		h.logger.Error("server::ApproveMaster::SaveMaster", err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusCreated)
	if _, err := rw.Write([]byte(fmt.Sprintf(`{ "id" : "%s" }`, masterID))); err != nil {
		h.logger.Error("server::ApproveMaster::Write", err)
		return
	}
	h.logger.Info("Response sent")
}
