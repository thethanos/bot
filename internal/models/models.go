package models

import (
	"gorm.io/gorm"
)

type City struct {
	gorm.Model
	Name string `gorm:"name"`
}

type ServiceCategory struct {
	gorm.Model
	Name string `gorm:"name"`
}

type Service struct {
	gorm.Model
	Name    string `gorm:"name"`
	CatID   uint   `gorm:"cat_id"`
	CatName string `gorm:"cat_name"`
}

type MasterServRelation struct {
	gorm.Model
	MasterID    uint   `gorm:"master_id"`
	Name        string `gorm:"name"`
	Description string `gorm:"description"`
	Contact     string `gorm:"contact"`
	CityID      uint   `gorm:"city_id"`
	CityName    string `gorm:"city_name"`
	ServCatID   uint   `gorm:"serv_cat_id"`
	ServCatName string `gorm:"serv_cat_name"`
	ServID      uint   `gorm:"serv_id"`
	ServName    string `gorm:"serv_name"`
}

type MasterRegForm struct {
	gorm.Model
	MasterID    uint   `gorm:"master_id"`
	Name        string `gorm:"name"`
	Description string `gorm:"description"`
	Contact     string `gorm:"contact"`
	CityID      uint   `gorm:"city_id"`
	CityName    string `gorm:"city_name"`
	ServCatID   uint   `gorm:"serv_cat_id"`
	ServCatName string `gorm:"serv_cat_name"`
	ServID      uint   `gorm:"serv_id"`
	ServName    string `gorm:"serv_name"`
}

type MasterImages struct {
	ID       uint   `gorm:"primaryKey;autoIncrement;notNull"`
	MasterID uint   `gorm:"master_id"`
	URL      string `gorm:"url"`
}
