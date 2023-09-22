package mapper

import (
	"bot/internal/entities"
	"bot/internal/models"
)

func FromCityModel(model *models.City) *entities.City {
	return &entities.City{
		ID:   model.ID,
		Name: model.Name,
	}
}

func FromServCatModel(model *models.ServiceCategory) *entities.ServiceCategory {
	return &entities.ServiceCategory{
		ID:   model.ID,
		Name: model.Name,
	}
}

func FromServiceModel(model *models.Service) *entities.Service {
	return &entities.Service{
		ID:      model.ID,
		Name:    model.Name,
		CatID:   model.CatID,
		CatName: model.CatName,
	}
}
