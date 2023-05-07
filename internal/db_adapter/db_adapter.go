package db_adapter

import (
	"database/sql"

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
	if err := d.dbConn.AutoMigrate(&City{}); err != nil {
		return err
	}
	if err := d.dbConn.AutoMigrate(&Service{}); err != nil {
		return err
	}
	if err := d.dbConn.AutoMigrate(&Master{}); err != nil {
		return err
	}

	return nil
}

func (d *DbAdapter) GetCities(service *Service) ([]*City, error) {
	cities := make([]*City, 0)
	if service == nil {
		tx := d.dbConn.Find(&cities)
		return cities, tx.Error
	}

	tx := d.dbConn.Where("service_id == ?", service.ID).Find(&cities)
	return cities, tx.Error
}

func (d *DbAdapter) GetServices(city *City) ([]*Service, error) {
	services := make([]*Service, 0)
	if city == nil {
		tx := d.dbConn.Find(&services)
		return services, tx.Error
	}

	tx := d.dbConn.Where("city_id == ?", city.ID).Find(&services)
	return services, tx.Error
}

func (d *DbAdapter) GetMasters(city *City, service *Service) ([]*Master, error) {

	masters := make([]*Master, 0)
	if city == nil && service == nil {
		tx := d.dbConn.Find(&masters)
		return masters, tx.Error
	}

	if city != nil && service == nil {
		tx := d.dbConn.Where("city_id == ?", city.ID).Find(&masters)
		return masters, tx.Error
	}

	if service != nil && city == nil {
		tx := d.dbConn.Where("service_id == ?", service.ID).Find(&masters)
		return masters, tx.Error
	}

	d.dbConn.Where("city_id == ? AND service_id == ?", city.ID, service.ID).Find(&masters)
	return masters, nil
}

func (d *DbAdapter) Test() {

	cities := []City{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name:      "Tel-Aviv",
			ServiceID: 1,
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Name:      "Jerusalem",
			ServiceID: 2,
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			Name:      "Netanya",
			ServiceID: 3,
		},
	}

	for _, city := range cities {
		d.dbConn.Create(&city)
	}

	masters := []Master{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name:      "Masha",
			CityID:    1,
			ServiceID: 1,
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Name:      "Sasha",
			CityID:    2,
			ServiceID: 2,
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			Name:      "Pasha",
			CityID:    3,
			ServiceID: 3,
		},
	}

	for _, master := range masters {
		d.dbConn.Create(&master)
	}

	services := []Service{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name:     "Service1",
			MasterID: 1,
			CityID:   1,
		},
		{
			Model: gorm.Model{
				ID: 2,
			},
			Name:     "Service2",
			CityID:   2,
			MasterID: 2,
		},
		{
			Model: gorm.Model{
				ID: 3,
			},
			Name:     "Service3",
			MasterID: 3,
			CityID:   3,
		},
	}

	for _, service := range services {
		d.dbConn.Create(&service)
	}
}
