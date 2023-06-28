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
		Image1:      model.Image1,
		Image2:      model.Image2,
		Image3:      model.Image3,
		Description: model.Description,
		Contact:     model.Contact,
		CityID:      model.CityID,
	}
}
