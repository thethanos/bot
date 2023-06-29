package db_adapter

import (
	"fmt"
	"multimessenger_bot/internal/config"
	"multimessenger_bot/internal/entities"
	"multimessenger_bot/internal/mapper"
	"multimessenger_bot/internal/models"
	"time"

	"multimessenger_bot/internal/logger"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DbAdapter struct {
	logger logger.Logger
	dbConn *gorm.DB
}

func NewDbAdapter(logger logger.Logger, cfg *config.Config) (*DbAdapter, error) {

	psqlconf := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PsqlHost,
		cfg.PsqlPort,
		cfg.PsqlUser,
		cfg.PsqlPass,
		cfg.PsqlDb,
	)

	dbConn, err := gorm.Open(postgres.Open(psqlconf), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DbAdapter{logger: logger, dbConn: dbConn}, nil
}

func (d *DbAdapter) AutoMigrate() error {
	if err := d.dbConn.AutoMigrate(&models.City{}); err != nil {
		return err
	}
	if err := d.dbConn.AutoMigrate(&models.ServiceCategory{}); err != nil {
		return err
	}
	if err := d.dbConn.AutoMigrate(&models.Service{}); err != nil {
		return err
	}
	if err := d.dbConn.AutoMigrate(&models.Master{}); err != nil {
		return err
	}
	if err := d.dbConn.AutoMigrate(&models.MasterRegForm{}); err != nil {
		return err
	}
	if err := d.dbConn.AutoMigrate(&models.Join{}); err != nil {
		return err
	}
	if err := d.dbConn.AutoMigrate(&models.JoinCityCategory{}); err != nil {
		return err
	}
	d.logger.Info("Auto-migration: success")
	return nil
}

func (d *DbAdapter) GetCity(name string) (*entities.City, error) {
	city := &models.City{}
	tx := d.dbConn.Where("index_str = ?", clearText(name)).First(city)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return mapper.FromCityModel(city), nil
}

func (d *DbAdapter) GetCities(serviceId string) ([]*entities.City, error) {

	result := make([]*entities.City, 0)
	cities := make([]*models.City, 0)

	if serviceId == "" {
		if err := d.dbConn.Find(&cities).Error; err != nil {
			d.logger.Error(err)
			return nil, err
		}
		for _, city := range cities {
			result = append(result, &entities.City{ID: city.ID, IndexStr: city.IndexStr, Name: city.Name})
		}
		return result, nil
	}

	joins := make([]*models.Join, 0)
	if err := d.dbConn.Where("service_id = ?", serviceId).Find(&joins).Error; err != nil {
		d.logger.Error(err)
		return nil, err
	}

	cityIds := make([]string, 0)
	for _, join := range joins {
		cityIds = append(cityIds, join.CityID)
	}

	if err := d.dbConn.Where("id IN ?", cityIds).Find(&cities).Error; err != nil {
		d.logger.Error(err)
		return nil, err
	}
	for _, city := range cities {
		result = append(result, &entities.City{ID: city.ID, Name: city.Name})
	}
	return result, nil
}

func (d *DbAdapter) GetCategories(cityId string) ([]*entities.ServiceCategory, error) {
	result := make([]*entities.ServiceCategory, 0)
	categories := make([]*models.ServiceCategory, 0)

	if len(cityId) == 0 {
		if err := d.dbConn.Find(&categories).Error; err != nil {
			d.logger.Error(err)
			return nil, err
		}
		for _, category := range categories {
			result = append(result, &entities.ServiceCategory{ID: category.ID, Name: category.Name})
		}
		return result, nil
	}

	joins := make([]*models.JoinCityCategory, 0)
	if err := d.dbConn.Where("city_id = ?", cityId).Find(&joins).Error; err != nil {
		return nil, err
	}

	categoryIds := make([]string, 0)
	for _, join := range joins {
		categoryIds = append(categoryIds, join.ServiceCategoryID)
	}

	if err := d.dbConn.Where("id IN ?", categoryIds).Find(&categories).Error; err != nil {
		return nil, err
	}
	for _, category := range categories {
		result = append(result, &entities.ServiceCategory{ID: category.ID, Name: category.Name})
	}
	return result, nil
}

func (d *DbAdapter) GetServices(categoryId, cityId string) ([]*entities.Service, error) {
	result := make([]*entities.Service, 0)
	services := make([]*models.Service, 0)

	if len(categoryId) == 0 {
		if err := d.dbConn.Find(&services).Error; err != nil {
			return nil, err
		}
	} else if len(cityId) != 0 {
		joins := make([]*models.Join, 0)
		if err := d.dbConn.Select("service_id").Distinct().Where("city_id = ?", cityId).Find(&joins).Error; err != nil {
			return nil, err
		}

		serviceIds := make([]string, 0)
		for _, join := range joins {
			serviceIds = append(serviceIds, join.ServiceID)
		}

		if err := d.dbConn.Where("category_id = ? AND id IN ?", categoryId, serviceIds).Find(&services).Error; err != nil {
			return nil, err
		}
	} else {
		if err := d.dbConn.Where("category_id = ?", categoryId).Find(&services).Error; err != nil {
			return nil, err
		}
	}

	for _, service := range services {
		result = append(result, &entities.Service{ID: service.ID, Name: service.Name, CategoryID: service.CategoryID})
	}
	return result, nil
}

/*
	func (d *DbAdapter) GetMasterPreview(id string) (*entities.Master, error) {
		master := &models.MasterPreview{}
		tx := d.dbConn.Where("id = ?", id).First(master)
		if tx.Error != nil {
			return nil, tx.Error
		}
		return &entities.Master{
			Name:        master.Name,
			Images:      master.Images,
			Description: master.Description,
		}, nil
	}
*/
func (d *DbAdapter) GetMasters(cityId, serviceId string) ([]*entities.Master, error) {
	result := make([]*entities.Master, 0)
	masters := make([]*models.Master, 0)
	joins := make([]*models.Join, 0)
	if err := d.dbConn.Where("city_id = ? AND service_id = ?", cityId, serviceId).Find(&joins).Error; err != nil {
		return nil, err
	}

	masterIds := make([]string, 0)
	for _, join := range joins {
		masterIds = append(masterIds, join.MasterID)
	}

	if err := d.dbConn.Where("id IN ?", masterIds).Find(&masters).Error; err != nil {
		return nil, err
	}
	for _, master := range masters {
		result = append(result, &entities.Master{
			ID:          master.ID,
			Name:        master.Name,
			Images:      master.Images,
			Description: master.Description,
			CityID:      master.CityID,
		})
	}
	return result, nil
}

func (d *DbAdapter) GetMasterRegForm(master_id string) (*entities.MasterRegForm, error) {

	master := &models.MasterRegForm{}
	if err := d.dbConn.Where("id = ?", master_id).First(master).Error; err != nil {
		return nil, err
	}

	return mapper.FromMasterRegFormModel(master), nil
}

func (d *DbAdapter) SaveServiceCategory(name string) error {
	id := fmt.Sprintf("%d", time.Now().Unix())
	service := &models.ServiceCategory{
		ID:   id,
		Name: name,
	}
	if err := d.dbConn.Create(service).Error; err != nil {
		return err
	}
	d.logger.Infof("New service category added successfully, id: %s, name: %s", id, name)
	return nil
}

func (d *DbAdapter) SaveService(name, categoryId string) error {
	id := fmt.Sprintf("%d", time.Now().Unix())
	service := &models.Service{
		ID:         id,
		Name:       name,
		CategoryID: categoryId,
	}
	if err := d.dbConn.Create(service).Error; err != nil {
		return err
	}
	d.logger.Infof("New service added successfully, id: %s, name: %s", id, name)
	return nil
}

func (d *DbAdapter) SaveCity(name string) error {
	id := fmt.Sprintf("%d", time.Now().Unix())
	city := &models.City{
		ID:       id,
		IndexStr: clearText(name),
		Name:     name,
	}
	if err := d.dbConn.Create(city).Error; err != nil {
		return err
	}
	d.logger.Infof("New city added successfully, id: %s, name: %s", id, name)
	return nil
}

func (d *DbAdapter) SaveMaster(data *entities.MasterRegForm) error {

	master := &models.Master{
		ID:          data.ID,
		Name:        data.Name,
		Description: data.Description,
		CityID:      data.CityID,
		Images:      data.Images,
	}

	tx := d.dbConn.Begin()
	defer tx.Rollback()

	if err := tx.Create(master).Error; err != nil {
		return err
	}

	if err := tx.Where("city_id = ? AND service_category_id = ?", data.CityID, data.CategoryID).First(&models.JoinCityCategory{}).Error; err != nil {
		d.logger.Infof("Creating new join record - city_id: %s, service_category_id: %s", data.CityID, data.CategoryID)
		if err := tx.Create(&models.JoinCityCategory{CityID: data.CityID, ServiceCategoryID: data.CategoryID}).Error; err != nil {
			return err
		}
	}

	for _, serviceID := range data.ServiceIDs {
		if err := tx.Create(&models.Join{CityID: data.CityID, ServiceID: serviceID, MasterID: master.ID}).Error; err != nil {
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	d.logger.Infof("New master added successfully, id: %s, name: %s", master.ID, master.Name)
	return nil
}

func (d *DbAdapter) SaveMasterRegForm(master *entities.MasterRegForm) (string, error) {
	id := fmt.Sprintf("%d", time.Now().Unix())

	images := make([]string, 0)
	for _, image := range master.Images {
		images = append(images, fmt.Sprintf("https://bot-dev-domain.com/pages/images/%s/%s", id, image))
	}

	regForm := &models.MasterRegForm{
		ID:          id,
		Name:        master.Name,
		CityID:      master.CityID,
		CategoryID:  master.CategoryID,
		ServiceIDs:  master.ServiceIDs,
		Contact:     master.Contact,
		Description: master.Description,
		Images:      images,
	}
	if err := d.dbConn.Create(regForm).Error; err != nil {
		return "", err
	}
	d.logger.Infof("Form saved successfully, id: %s, name: %s", id, master.Name)
	return id, nil
}
