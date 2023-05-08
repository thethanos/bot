package models

type City struct {
	ID   string `gorm:"primarykey"`
	Name string `gorm:"name"`
}

type Service struct {
	ID   string `gorm:"primarykey"`
	Name string `gorm:"name"`
}

type Master struct {
	ID     string `gorm:"primarykey"`
	Name   string `gorm:"name"`
	CityID string `gorm:"city_id"`
}

type Join struct {
	CityID    string `gorm:"city_id"`
	ServiceID string `gorm:"service_id"`
	MasterID  string `gorm:"master_id"`
}
