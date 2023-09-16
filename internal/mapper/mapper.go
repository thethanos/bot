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

func FromMasterServRelationModel(model *models.MasterServRelation, imgURLs []string) *entities.Master {
	return &entities.Master{
		ID:          model.MasterID,
		Name:        model.Name,
		Description: model.Description,
		Contact:     model.Contact,
		CityName:    model.CityName,
		ServCatName: model.ServCatName,
		RegDate:     model.CreatedAt.Format("2006-01-02"),
		ImageURLs:   imgURLs,
	}
}
