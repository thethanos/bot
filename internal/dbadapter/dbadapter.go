package dbadapter

import (
	"bot/internal/config"
	"bot/internal/entities"
	"bot/internal/mapper"
	"bot/internal/models"
	"fmt"

	"bot/internal/logger"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBAdapter struct {
	logger logger.Logger
	cfg    *config.Config
	DBConn *gorm.DB
}

func NewDbAdapter(logger logger.Logger, cfg *config.Config) (*DBAdapter, error) {

	psqlconf := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.PsqlHost,
		cfg.PsqlPort,
		cfg.PsqlUser,
		cfg.PsqlPass,
		cfg.PsqlDb,
	)

	DBConn, err := gorm.Open(postgres.Open(psqlconf), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return &DBAdapter{logger: logger, cfg: cfg, DBConn: DBConn}, nil
}

func (d *DBAdapter) GetCities(servID uint, page, limit int) ([]*entities.City, error) {

	if servID != 0 {
		return d.GetCitiesByService(servID, page, limit)
	}

	cities := make([]*models.City, 0)
	if err := d.DBConn.Offset(page * limit).Limit(limit).Order("name ASC").Find(&cities).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.City, 0)
	for _, city := range cities {
		result = append(result, mapper.FromCityModel(city))
	}

	return result, nil
}

func (d *DBAdapter) GetCitiesByService(servID uint, page, limit int) ([]*entities.City, error) {

	relations := make([]*models.MasterServRelation, 0)
	subquery := d.DBConn.Table("master_serv_relations").Offset(page * limit).Limit(limit)
	subquery = subquery.Where("serv_id = ?", servID).Select("DISTINCT ON (city_id) city_id, city_name")
	if err := d.DBConn.Table("(?) as subquery", subquery).Order("city_name ASC").Find(&relations).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.City, 0)
	for _, relation := range relations {
		result = append(result, &entities.City{
			ID:   relation.CityID,
			Name: relation.CityName,
		})
	}

	return result, nil
}

func (d *DBAdapter) GetServCategories(cityID uint, page, limit int) ([]*entities.ServiceCategory, error) {

	if cityID != 0 {
		return d.GetServCategoriesByCity(cityID, page, limit)
	}

	categories := make([]*models.ServiceCategory, 0)
	if err := d.DBConn.Offset(page * limit).Limit(limit).Order("name ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.ServiceCategory, 0)
	for _, category := range categories {
		result = append(result, mapper.FromServCatModel(category))
	}

	return result, nil
}

func (d *DBAdapter) GetServCategoriesByCity(cityID uint, page, limit int) ([]*entities.ServiceCategory, error) {

	relations := make([]*models.MasterServRelation, 0)
	subquery := d.DBConn.Table("master_serv_relations").Offset(page * limit).Limit(limit)
	subquery = subquery.Where("city_id = ?", cityID).Select("DISTINCT ON (serv_cat_id) serv_cat_id, serv_cat_name")
	if err := d.DBConn.Table("(?) as subquery", subquery).Order("serv_cat_name ASC").Find(&relations).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.ServiceCategory, 0)
	for _, relation := range relations {
		result = append(result, &entities.ServiceCategory{
			ID:   relation.ServCatID,
			Name: relation.ServCatName,
		})
	}

	return result, nil
}

func (d *DBAdapter) GetServices(categoryID, cityID uint, page, limit int) ([]*entities.Service, error) {

	if cityID != 0 {
		return d.GetServicesByCity(categoryID, cityID, page, limit)
	}

	return d.GetServicesByCategory(categoryID, page, limit)
}

func (d *DBAdapter) GetServicesByCity(categoryID, cityID uint, page, limit int) ([]*entities.Service, error) {

	relations := make([]*models.MasterServRelation, 0)
	subquery := d.DBConn.Offset(page * limit).Limit(limit)
	if categoryID != 0 {
		subquery = subquery.Where("serv_cat_id = ?", categoryID)
	}

	subquery = subquery.Where("city_id = ?", cityID).Select("DISTINCT ON (serv_id) serv_id, serv_name, serv_cat_id, serv_cat_name")
	if err := d.DBConn.Table("(?) as subquery", subquery).Order("serv_name ASC").Find(&relations).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.Service, 0)
	for _, relation := range relations {
		result = append(result, &entities.Service{
			ID:      relation.ServID,
			Name:    relation.ServName,
			CatID:   relation.ServCatID,
			CatName: relation.ServCatName,
		})
	}

	return result, nil
}

func (d *DBAdapter) GetServicesByCategory(categoryID uint, page, limit int) ([]*entities.Service, error) {

	query := d.DBConn.Offset(page * limit).Limit(limit).Order("name ASC")
	if categoryID != 0 {
		query = query.Where("cat_id = ?", categoryID)
	}

	services := make([]*models.Service, 0)
	if err := query.Find(&services).Error; err != nil {
		return nil, err
	}

	result := make([]*entities.Service, 0)
	for _, service := range services {
		result = append(result, mapper.FromServiceModel(service))
	}

	return result, nil
}
