package dbadapter

import "bot/internal/entities"

type DBInterface interface {
	GetCities(servID string, page, limit int) ([]*entities.City, error)
	GetCitiesByService(servID string, page, limit int) ([]*entities.City, error)
	GetServCategories(cityID string, page, limit int) ([]*entities.ServiceCategory, error)
	GetServCategoriesByCity(cityID string, page, limit int) ([]*entities.ServiceCategory, error)
	GetServices(categoryID, cityID string, page, limit int) ([]*entities.Service, error)
	GetServicesByCity(categoryID, cityID string, page, limit int) ([]*entities.Service, error)
	GetServicesByCategory(categoryID string, page, limit int) ([]*entities.Service, error)
}
