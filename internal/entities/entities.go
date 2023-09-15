package entities

type UserState struct {
	City            *City
	ServiceCategory *ServiceCategory
	Service         *Service
	RawInput        map[string]string
}

func (u *UserState) Reset() {
	u.City = nil
	u.ServiceCategory = nil
	u.Service = nil
	u.RawInput = make(map[string]string)
}

func (u UserState) GetCityID() uint {
	if u.City != nil {
		return u.City.ID
	}
	return 0
}

func (u UserState) GetServiceID() uint {
	if u.Service != nil {
		return u.Service.ID
	}
	return 0
}

type City struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type ServiceCategory struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	ID      uint   `json:"id"`
	Name    string `json:"name" validate:"required"`
	CatID   uint   `json:"category_id" validate:"required"`
	CatName string `json:"category_name" validate:"required"`
}

type Master struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Contact     string   `json:"contact"`
	Images      []string `json:"images"`
	CityName    string   `json:"cityName"`
	ServCatName string   `json:"servCatName"`
}

type MasterRegForm struct {
	ID          uint     `json:"id,omitempty"`
	Name        string   `json:"name" validate:"required"`
	Images      []string `json:"images" validate:"required"`
	Description string   `json:"description,omitempty"`
	Contact     string   `json:"contact" validate:"required"`
	CityID      uint     `json:"city_id" validate:"required"`
	ServCatID   uint     `json:"serv_cat_id" validate:"required"`
	ServIDs     []string `json:"serv_ids" validate:"required"`
}
