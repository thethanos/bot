package dbadapter

import "bot/internal/entities"

type DBAdapterMock struct {
}

func (db *DBAdapterMock) GetCities(servID string, page, limit int) ([]*entities.City, error) {

	cities := []*entities.City{
		{ID: "0", Name: "Tel-Aviv"},
		{ID: "1", Name: "Jerusalem"},
		{ID: "2", Name: "Haifa"},
	}

	return cities, nil
}

func (db *DBAdapterMock) GetCitiesByService(servID string, page, limit int) ([]*entities.City, error) {

	cities := []*entities.City{
		{ID: "0", Name: "Tel-Aviv"},
		{ID: "1", Name: "Jerusalem"},
		{ID: "2", Name: "Haifa"},
	}

	return cities, nil
}

func (db *DBAdapterMock) GetServCategories(cityID string, page, limit int) ([]*entities.ServiceCategory, error) {

	categories := []*entities.ServiceCategory{
		{ID: "0", Name: "Face"},
		{ID: "1", Name: "Body"},
		{ID: "2", Name: "Head"},
	}

	return categories, nil
}

func (db *DBAdapterMock) GetServCategoriesByCity(cityID string, page, limit int) ([]*entities.ServiceCategory, error) {

	categories := []*entities.ServiceCategory{
		{ID: "0", Name: "Face"},
		{ID: "1", Name: "Body"},
		{ID: "2", Name: "Head"},
	}

	return categories, nil
}

func (db *DBAdapterMock) GetServices(categoryID, cityID string, page, limit int) ([]*entities.Service, error) {

	services := []*entities.Service{
		{ID: "0", Name: "Makeup", CatID: "0", CatName: "Face"},
		{ID: "1", Name: "Massage", CatID: "1", CatName: "Body"},
		{ID: "2", Name: "Haircut", CatID: "2", CatName: "Head"},
	}

	return services, nil
}

func (db *DBAdapterMock) GetServicesByCity(categoryID, cityID string, page, limit int) ([]*entities.Service, error) {

	services := []*entities.Service{
		{ID: "0", Name: "Makeup", CatID: "0", CatName: "Face"},
		{ID: "1", Name: "Massage", CatID: "1", CatName: "Body"},
		{ID: "2", Name: "Haircut", CatID: "2", CatName: "Head"},
	}

	return services, nil
}

func (db *DBAdapterMock) GetServicesByCategory(categoryID string, page, limit int) ([]*entities.Service, error) {

	services := []*entities.Service{
		{ID: "0", Name: "Makeup", CatID: "0", CatName: "Face"},
		{ID: "1", Name: "Massage", CatID: "1", CatName: "Body"},
		{ID: "2", Name: "Haircut", CatID: "2", CatName: "Head"},
	}

	return services, nil
}
