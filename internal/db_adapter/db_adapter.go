package db_adapter

import (
	"database/sql"
	"fmt"
	"multimessenger_bot/internal/entities"
	"multimessenger_bot/internal/models"
	"time"

	"go.mau.fi/whatsmeow/store/sqlstore"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DbAdapter struct {
	dbConn *gorm.DB
}

func NewDbAdapter() (*DbAdapter, *sqlstore.Container, error) {

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

	return &DbAdapter{dbConn: dbConn}, container, nil
}

func (d *DbAdapter) AutoMigrate() error {
	if err := d.dbConn.AutoMigrate(&models.City{}); err != nil {
		return err
	}
	if err := d.dbConn.AutoMigrate(&models.Service{}); err != nil {
		return err
	}
	if err := d.dbConn.AutoMigrate(&models.Master{}); err != nil {
		return err
	}
	if err := d.dbConn.AutoMigrate(&models.Join{}); err != nil {
		return err
	}
	return nil
}

func (d *DbAdapter) GetCities(serviceId string) ([]*entities.City, error) {

	result := make([]*entities.City, 0)
	cities := make([]*models.City, 0)

	if serviceId == "" {
		if err := d.dbConn.Find(&cities).Error; err != nil {
			return nil, err
		}
		for _, city := range cities {
			result = append(result, &entities.City{ID: city.ID, Name: city.Name})
		}
		return result, nil
	}

	joins := make([]*models.Join, 0)
	if err := d.dbConn.Where("service_id == ?", serviceId).Find(&joins).Error; err != nil {
		return nil, err
	}

	cityIds := make([]string, 0)
	for _, join := range joins {
		cityIds = append(cityIds, join.CityID)
	}

	if err := d.dbConn.Where("id IN ?", cityIds).Find(&cities).Error; err != nil {
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

	tx := d.dbConn.Create(&models.Join{CityID: data.City.ID, ServiceID: data.Service.ID, MasterID: id})

	return tx.Error
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

func (d *DbAdapter) Test() error {

	cities := []*models.City{
		{
			ID:   "1",
			Name: "Тель-Авив",
		},
		{
			ID:   "2",
			Name: "Хайфа",
		},
		{
			ID:   "3",
			Name: "Иерусалим",
		},
		{
			ID:   "4",
			Name: "Нетания",
		},
	}
	/*
		masters := []*Master{
			{
				Model: gorm.Model{
					ID: 1,
				},
				Name: "Наталья",
			},
			{
				Model: gorm.Model{
					ID: 2,
				},
				Name: "Мария",
			},
			{
				Model: gorm.Model{
					ID: 3,
				},
				Name: "Александра",
			},
			{
				Model: gorm.Model{
					ID: 4,
				},
				Name: "Юлия",
			},
		}
	*/
	services := []*models.Service{
		{
			ID:   "1",
			Name: "Наращивание ресниц",
		},
		{
			ID:   "2",
			Name: "Окрашивание бровей",
		},
		{
			ID:   "3",
			Name: "Окрашивание ресниц",
		},
		{
			ID:   "4",
			Name: "Снятие ресниц",
		},
	}

	for _, service := range services {
		if err := d.dbConn.Create(service).Error; err != nil {
			fmt.Println(err)
			return err
		}
	}

	for _, city := range cities {
		if err := d.dbConn.Create(city).Error; err != nil {
			fmt.Println(err)
			return err
		}
	}

	return nil
}
