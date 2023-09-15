package models

import (
	"database/sql/driver"
	"errors"
	"strings"

	"gorm.io/gorm"
)

type string_array []string

func (i *string_array) Scan(src any) error {
	bytes, ok := src.(string)
	if !ok {
		return errors.New("src value cannot cast to []byte")
	}
	*i = strings.Split(string(bytes), ",")
	return nil
}
func (i string_array) Value() (driver.Value, error) {
	if len(i) == 0 {
		return nil, nil
	}
	return strings.Join(i, ","), nil
}

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
	MasterID    uint         `gorm:"master_id"`
	Name        string       `gorm:"name"`
	Description string       `gorm:"description"`
	Contact     string       `gorm:"contact"`
	Images      string_array `gorm:"type:text"`
	CityID      uint         `gorm:"city_id"`
	CityName    string       `gorm:"city_name"`
	ServCatID   uint         `gorm:"serv_cat_id"`
	ServCatName string       `gorm:"serv_cat_name"`
	ServID      uint         `gorm:"serv_id"`
	ServName    string       `gorm:"serv_name"`
}

type MasterRegForm struct {
	gorm.Model
	Name        string       `gorm:"name"`
	CityID      uint         `gorm:"city_id"`
	ServCatID   uint         `gorm:"service_category_id"`
	ServIDs     string_array `gorm:"type:text"`
	Contact     string       `gorm:"contact"`
	Images      string_array `gorm:"type:text"`
	Description string       `gotm:"description"`
}
