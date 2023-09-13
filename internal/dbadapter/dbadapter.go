package dbadapter

import (
	"bot/internal/config"
	"bot/internal/entities"
	"bot/internal/mapper"
	"bot/internal/models"
	"fmt"
	"time"

	"bot/internal/logger"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBAdapter struct {
	logger logger.Logger
	cfg    *config.Config
	DBConn *gorm.DB
}

func NewDbAdapter(logger logger.Logger, cfg *config.Config) (*DBAdapter, error) {

	psqlconf := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PsqlHost,
		cfg.PsqlPort,
		cfg.PsqlUser,
		cfg.PsqlPass,
		cfg.PsqlDb,
	)

	DBConn, err := gorm.Open(postgres.Open(psqlconf), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DBAdapter{logger: logger, cfg: cfg, DBConn: DBConn}, nil
}

func (d *DBAdapter) AutoMigrate() error {
	if err := d.DBConn.AutoMigrate(&models.City{}); err != nil {
		return err
	}
	if err := d.DBConn.AutoMigrate(&models.ServiceCategory{}); err != nil {
		return err
	}
	if err := d.DBConn.AutoMigrate(&models.Service{}); err != nil {
		return err
	}
	if err := d.DBConn.AutoMigrate(&models.Master{}); err != nil {
		return err
	}
	if err := d.DBConn.AutoMigrate(&models.MasterRegForm{}); err != nil {
		return err
	}
	if err := d.DBConn.AutoMigrate(&models.Join{}); err != nil {
		return err
	}
	if err := d.DBConn.AutoMigrate(&models.JoinCityCategory{}); err != nil {
		return err
	}
	d.logger.Info("Auto-migration: success")
	return nil
}

func (d *DBAdapter) GetCities(serviceId string, page, limit int) ([]*entities.City, error) {

	result := make([]*entities.City, 0)
	cities := make([]*models.City, 0)

	if serviceId == "" {
		if err := d.DBConn.Offset(page).Limit(limit).Find(&cities).Error; err != nil {
			d.logger.Error(err)
			return nil, err
		}
		for _, city := range cities {
			result = append(result, &entities.City{ID: city.ID, Name: city.Name})
		}
		return result, nil
	}

	joins := make([]*models.Join, 0)
	if err := d.DBConn.Where("service_id = ?", serviceId).Find(&joins).Error; err != nil {
		d.logger.Error(err)
		return nil, err
	}

	cityIds := make([]string, 0)
	for _, join := range joins {
		cityIds = append(cityIds, join.CityID)
	}

	if err := d.DBConn.Offset(page).Limit(limit).Where("id IN ?", cityIds).Find(&cities).Error; err != nil {
		d.logger.Error(err)
		return nil, err
	}
	for _, city := range cities {
		result = append(result, &entities.City{ID: city.ID, Name: city.Name})
	}
	return result, nil
}

func (d *DBAdapter) GetServiceCategories(cityId string, page, limit int) ([]*entities.ServiceCategory, error) {
	result := make([]*entities.ServiceCategory, 0)
	categories := make([]*models.ServiceCategory, 0)

	if len(cityId) == 0 {
		if err := d.DBConn.Offset(page).Limit(limit).Find(&categories).Error; err != nil {
			d.logger.Error(err)
			return nil, err
		}
		for _, category := range categories {
			result = append(result, &entities.ServiceCategory{ID: category.ID, Name: category.Name})
		}
		return result, nil
	}

	joins := make([]*models.JoinCityCategory, 0)
	if err := d.DBConn.Where("city_id = ?", cityId).Find(&joins).Error; err != nil {
		return nil, err
	}

	categoryIds := make([]string, 0)
	for _, join := range joins {
		categoryIds = append(categoryIds, join.ServiceCategoryID)
	}

	if err := d.DBConn.Offset(page).Limit(limit).Where("id IN ?", categoryIds).Find(&categories).Error; err != nil {
		return nil, err
	}
	for _, category := range categories {
		result = append(result, &entities.ServiceCategory{ID: category.ID, Name: category.Name})
	}
	return result, nil
}

func (d *DBAdapter) GetServices(categoryId, cityId string, page, limit int) ([]*entities.Service, error) {
	result := make([]*entities.Service, 0)
	services := make([]*models.Service, 0)

	if len(categoryId) == 0 {
		if err := d.DBConn.Offset(page).Limit(limit).Find(&services).Error; err != nil {
			return nil, err
		}
	} else if len(cityId) != 0 {
		joins := make([]*models.Join, 0)
		if err := d.DBConn.Select("service_id").Distinct().Where("city_id = ?", cityId).Find(&joins).Error; err != nil {
			return nil, err
		}

		serviceIds := make([]string, 0)
		for _, join := range joins {
			serviceIds = append(serviceIds, join.ServiceID)
		}

		if err := d.DBConn.Offset(page).Limit(limit).Where("category_id = ? AND id IN ?", categoryId, serviceIds).Find(&services).Error; err != nil {
			return nil, err
		}
	} else {
		if err := d.DBConn.Offset(page).Limit(limit).Where("category_id = ?", categoryId).Find(&services).Error; err != nil {
			return nil, err
		}
	}

	for _, service := range services {
		result = append(result, &entities.Service{ID: service.ID, Name: service.Name, CategoryID: service.CategoryID})
	}
	return result, nil
}

func (d *DBAdapter) GetMasters(cityId, serviceId string, page, limit int) ([]*entities.Master, error) {
	result := make([]*entities.Master, 0)
	masters := make([]*models.Master, 0)
	joins := make([]*models.Join, 0)

	query := d.DBConn.Offset(page * limit).Limit(limit)
	if len(cityId) != 0 || len(serviceId) != 0 {
		if len(cityId) != 0 {
			query = query.Where("city_id = ?", cityId)
		}
		if len(serviceId) != 0 {
			query = query.Where("service_id = ?", serviceId)
		}

		if err := query.Find(&joins).Error; err != nil {
			return nil, err
		}

		masterIds := make([]string, 0)
		for _, join := range joins {
			masterIds = append(masterIds, join.MasterID)
		}

		if err := d.DBConn.Where("id IN ?", masterIds).Find(&masters).Error; err != nil {
			return nil, err
		}
	} else {
		if err := query.Find(&masters).Error; err != nil {
			return nil, err
		}
	}

	for _, master := range masters {
		result = append(result, &entities.Master{
			ID:          master.ID,
			Name:        master.Name,
			Images:      master.Images,
			Description: master.Description,
		})
	}
	return result, nil
}

func (d *DBAdapter) GetMasterRegForm(master_id string) (*entities.MasterRegForm, error) {

	master := &models.MasterRegForm{}
	if err := d.DBConn.Where("id = ?", master_id).First(master).Error; err != nil {
		return nil, err
	}

	return mapper.FromMasterRegFormModel(master), nil
}

func (d *DBAdapter) SaveServiceCategory(name string) (string, error) {
	id := fmt.Sprintf("%d", time.Now().Unix())
	service := &models.ServiceCategory{
		ID:   id,
		Name: name,
	}
	if err := d.DBConn.Create(service).Error; err != nil {
		return "", err
	}
	d.logger.Infof("New service category added successfully, id: %s, name: %s", id, name)
	return id, nil
}

func (d *DBAdapter) SaveService(name, categoryId string) (string, error) {
	id := fmt.Sprintf("%d", time.Now().Unix())
	service := &models.Service{
		ID:         id,
		Name:       name,
		CategoryID: categoryId,
	}
	if err := d.DBConn.Create(service).Error; err != nil {
		return "", err
	}
	d.logger.Infof("New service added successfully, id: %s, name: %s", id, name)
	return id, nil
}

func (d *DBAdapter) SaveCity(name string) (string, error) {
	id := fmt.Sprintf("%d", time.Now().Unix())
	city := &models.City{
		ID:   id,
		Name: name,
	}
	if err := d.DBConn.Create(city).Error; err != nil {
		return "", err
	}
	d.logger.Infof("New city added successfully, id: %s, name: %s", id, name)
	return id, nil
}

func (d *DBAdapter) SaveMaster(data *entities.MasterRegForm) (string, error) {

	master := &models.Master{
		ID:          data.ID,
		Name:        data.Name,
		Contact:     data.Contact,
		Description: data.Description,
		Images:      data.Images,
		CityID:      data.CityID,
	}

	tx := d.DBConn.Begin()
	defer tx.Rollback()

	if err := tx.Create(master).Error; err != nil {
		return "", err
	}

	if err := tx.Where("city_id = ? AND service_category_id = ?", data.CityID, data.ServiceCategoryID).First(&models.JoinCityCategory{}).Error; err != nil {
		d.logger.Infof("Creating new join record - city_id: %s, service_category_id: %s", data.CityID, data.ServiceCategoryID)
		if err := tx.Create(&models.JoinCityCategory{CityID: data.CityID, ServiceCategoryID: data.ServiceCategoryID}).Error; err != nil {
			return "", err
		}
	}

	for _, serviceID := range data.ServiceIDs {
		if err := tx.Create(&models.Join{CityID: data.CityID, ServiceID: serviceID, MasterID: master.ID}).Error; err != nil {
			return "", err
		}
	}

	if err := tx.Delete(&models.MasterRegForm{ID: data.ID}).Error; err != nil {
		return "", err
	}

	if err := tx.Commit().Error; err != nil {
		return "", err
	}

	d.logger.Infof("New master added successfully, id: %s, name: %s", master.ID, master.Name)
	return master.ID, nil
}

func (d *DBAdapter) SaveMasterRegForm(master *entities.MasterRegForm) (string, error) {
	id := fmt.Sprintf("%d", time.Now().Unix())

	images := make([]string, 0)
	for _, image := range master.Images {
		images = append(images, fmt.Sprintf("%s/%s/%s", d.cfg.ImagePrefix, id, image))
	}

	regForm := &models.MasterRegForm{
		ID:                id,
		Name:              master.Name,
		CityID:            master.CityID,
		ServiceCategoryID: master.ServiceCategoryID,
		ServiceIDs:        master.ServiceIDs,
		Contact:           master.Contact,
		Description:       master.Description,
		Images:            images,
	}
	if err := d.DBConn.Create(regForm).Error; err != nil {
		return "", err
	}
	d.logger.Infof("Form saved successfully, id: %s, name: %s", id, master.Name)
	return id, nil
}
