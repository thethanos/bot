package models

import (
	"database/sql/driver"
	"errors"
	"strings"
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
	ID       string `gorm:"primarykey"`
	IndexStr string `gorm:"index_str"`
	Name     string `gorm:"name"`
}

type ServiceCategory struct {
	ID       string `gorm:"primarykey"`
	IndexStr string `gorm:"index_str"`
	Name     string `gorm:"name"`
}

type Service struct {
	ID         string `gorm:"primarykey"`
	IndexStr   string `gorm:"index_str"`
	Name       string `gorm:"name"`
	CategoryID string `gorm:"category_id"`
}

type Master struct {
	ID          string       `gorm:"primarykey"`
	IndexStr    string       `gorm:"index_str"`
	Name        string       `gorm:"name"`
	Images      string_array `gorm:"type:text"`
	Description string       `gorm:"description"`
	Contact     string       `gorm:"contact"`
	CityID      string       `gorm:"city_id"`
}

type MasterRegForm struct {
	ID         string       `gorm:"primarykey"`
	Name       string       `gorm:"name"`
	CityID     string       `gorm:"city_id"`
	CategoryID string       `gorm:"category_id"`
	ServiceIDs string_array `gorm:"type:text"`
	Contact    string       `gorm:"contact"`
}

type Join struct {
	CityID    string `gorm:"city_id"`
	ServiceID string `gorm:"service_id"`
	MasterID  string `gorm:"master_id"`
}

type JoinCityCategory struct {
	CityID            string `gorm:"city_id"`
	ServiceCategoryID string `gorm:"service_category_id"`
}
