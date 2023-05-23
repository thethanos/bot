package models

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
	ID          string `gorm:"primarykey"`
	IndexStr    string `gorm:"index_str"`
	Name        string `gorm:"name"`
	Image       string `gorm:"image"`
	Description string `gorm:"description"`
	CityID      string `gorm:"city_id"`
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
