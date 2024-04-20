package models

type City struct {
	ID   string `gorm:"column:id;type:varchar(36);"`
	Name string `gorm:"name"`
}

type ServiceCategory struct {
	ID   string `gorm:"column:id;type:varchar(36);"`
	Name string `gorm:"name"`
}

type Service struct {
	ID      string `gorm:"column:id;type:varchar(36);"`
	Name    string `gorm:"name"`
	CatID   string `gorm:"column:cat_id;type:varchar(36);"`
	CatName string `gorm:"cat_name"`
}

type MasterServRelation struct {
	ID          string `gorm:"column:id;type:varchar(36);"`
	MasterID    string `gorm:"column:master_id;type:varchar(36);"`
	Name        string `gorm:"name"`
	Description string `gorm:"description"`
	Contact     string `gorm:"contact"`
	CityID      string `gorm:"column:city_id;type:varchar(36);"`
	CityName    string `gorm:"city_name"`
	ServCatID   string `gorm:"column:serv_cat_id;type:varchar(36);"`
	ServCatName string `gorm:"serv_cat_name"`
	ServID      string `gorm:"column:serv_id;type:varchar(36);"`
	ServName    string `gorm:"serv_name"`
}
