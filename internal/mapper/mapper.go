package mapper

import (
	"multimessenger_bot/internal/entities"
	"multimessenger_bot/internal/models"
)

func FromCityModel(model *models.City) *entities.City {
	return &entities.City{
		ID:   model.ID,
		Name: model.Name,
	}
}

func FromServiceModel(model *models.Service) *entities.Service {
	return &entities.Service{
		ID:   model.ID,
		Name: model.Name,
	}
}

func FromMasterModel(model *models.Master) *entities.Master {
	return &entities.Master{
		ID:          model.ID,
		Name:        model.Name,
		Images:      model.Images,
		Description: model.Description,
		CityID:      model.CityID,
	}
}
