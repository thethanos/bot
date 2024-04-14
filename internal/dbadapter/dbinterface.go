package dbadapter

import "bot/internal/entities"

type DBInterface interface {
	GetCities(servID uint, page, limit int) ([]*entities.City, error)
	GetCitiesByService(servID uint, page, limit int) ([]*entities.City, error)
	GetServCategories(cityID uint, page, limit int) ([]*entities.ServiceCategory, error)
	GetServCategoriesByCity(cityID uint, page, limit int) ([]*entities.ServiceCategory, error)
	GetServices(categoryID, cityID uint, page, limit int) ([]*entities.Service, error)
	GetServicesByCity(categoryID, cityID uint, page, limit int) ([]*entities.Service, error)
	GetServicesByCategory(categoryID uint, page, limit int) ([]*entities.Service, error)
}
