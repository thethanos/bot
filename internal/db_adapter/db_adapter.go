package db_adapter

import (
	"database/sql"
	"fmt"
	"multimessenger_bot/internal/entities"
	"multimessenger_bot/internal/models"
	"time"

	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbAdapter struct {
	logger *zap.SugaredLogger
	dbConn *gorm.DB
}

func NewDbAdapter(logger *zap.SugaredLogger) (*DbAdapter, *sqlstore.Container, error) {

	rawDbConn, err := sql.Open("sqlite3", "file:sqlite.db?_foreign_keys=on")
	if err != nil {
		return nil, nil, err
	}

	dbConn, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	container := sqlstore.NewWithDB(rawDbConn, "sqlite3", nil)
	if err := container.Upgrade(); err != nil {
		return nil, nil, err
	}

	return &DbAdapter{logger: logger, dbConn: dbConn}, container, nil
}

func (d *DbAdapter) AutoMigrate() error {
	if err := d.dbConn.AutoMigrate(&models.City{}); err != nil {
		d.logger.Error("db_adapter::DbAdapter::AutoMigrate", err)
		return err
	}
	if err := d.dbConn.AutoMigrate(&models.Service{}); err != nil {
		d.logger.Error("db_adapter::DbAdapter::AutoMigrate", err)
		return err
	}
	if err := d.dbConn.AutoMigrate(&models.Master{}); err != nil {
		d.logger.Error("db_adapter::DbAdapter::AutoMigrate", err)
		return err
	}
	if err := d.dbConn.AutoMigrate(&models.Join{}); err != nil {
		d.logger.Error("db_adapter::DbAdapter::AutoMigrate", err)
		return err
	}

	d.logger.Info("Auto-migration: success")
	return nil
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
			result = append(result, &entities.City{ID: city.ID, Name: city.Name})
		}
		return result, nil
	}

	joins := make([]*models.Join, 0)
	if err := d.dbConn.Where("service_id == ?", serviceId).Find(&joins).Error; err != nil {
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

func (d *DbAdapter) GetServices(cityId string) ([]*entities.Service, error) {
	result := make([]*entities.Service, 0)
	services := make([]*models.Service, 0)

	if cityId == "" {
		if err := d.dbConn.Find(&services).Error; err != nil {
			d.logger.Error(err)
			return nil, err
		}
		for _, service := range services {
			result = append(result, &entities.Service{ID: service.ID, Name: service.Name})
		}
		return result, nil
	}

	joins := make([]*models.Join, 0)
	if err := d.dbConn.Where("city_id == ?", cityId).Find(&joins).Error; err != nil {
		return nil, err
	}

	serviceIds := make([]string, 0)
	for _, join := range joins {
		serviceIds = append(serviceIds, join.ServiceID)
	}

	if err := d.dbConn.Where("id IN ?", serviceIds).Find(&services).Error; err != nil {
		return nil, err
	}
	for _, service := range services {
		result = append(result, &entities.Service{ID: service.ID, Name: service.Name})
	}
	return result, nil
}

func (d *DbAdapter) GetMasters(cityId, serviceId string) ([]*entities.Master, error) {
	result := make([]*entities.Master, 0)
	masters := make([]*models.Master, 0)

	if cityId == "" && serviceId == "" {
		if err := d.dbConn.Find(&masters).Error; err != nil {
			return nil, err
		}
		for _, master := range masters {
			result = append(result, &entities.Master{ID: master.ID, Name: master.Name})
		}
		return result, nil
	}

	joins := make([]*models.Join, 0)
	if err := d.dbConn.Where("city_id == ? AND service_id == ?", cityId, serviceId).Find(&joins).Error; err != nil {
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
		result = append(result, &entities.Master{ID: master.ID, Name: master.Name})
	}
	return result, nil
}

func (d *DbAdapter) SaveNewCity(name string) error {
	id := fmt.Sprintf("%d", time.Now().Unix())
	city := &models.City{
		ID:   id,
		Name: name,
	}
	return d.dbConn.Create(city).Error
}

func (d *DbAdapter) SaveNewService(name string) error {
	id := fmt.Sprintf("%d", time.Now().Unix())
	service := &models.Service{
		ID:   id,
		Name: name,
	}
	return d.dbConn.Create(service).Error
}

func (d *DbAdapter) SaveNewMaster(data *entities.UserState) error {
	id := fmt.Sprintf("%d", time.Now().Unix())
	master := &models.Master{
		ID:     id,
		Name:   data.RawInput["name"],
		CityID: data.City.ID,
	}

	if err := d.dbConn.Create(master).Error; err != nil {
		return err
	}

	return d.dbConn.Create(&models.Join{CityID: data.City.ID, ServiceID: data.Service.ID, MasterID: id}).Error
}

func (d *DbAdapter) getCityByName(name string) (*models.City, error) {
	city := &models.City{}
	tx := d.dbConn.Where("name == ?", name).Find(city)
	return city, tx.Error
}

func (d *DbAdapter) getServiceByName(name string) (*models.Service, error) {
	service := &models.Service{}
	tx := d.dbConn.Where("name == ?", name).Find(service)
	return service, tx.Error
}
