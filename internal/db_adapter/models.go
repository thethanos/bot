package db_adapter

import "gorm.io/gorm"

type City struct {
	gorm.Model
	Name      string `gorm:"name"`
	ServiceID uint   `gorm:"service_id"`
}

type Service struct {
	gorm.Model
	Name     string `gorm:"name"`
	CityID   uint   `gorm:"city_id"`
	MasterID uint   `gorm:"master_id"`
}

type Master struct {
	gorm.Model
	Name      string `gorm:"name"`
	CityID    uint   `gorm:"city_id"`
	ServiceID uint   `gorm:"service_id"`
}
