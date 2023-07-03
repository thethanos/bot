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

func (u UserState) GetCityID() string {
	if u.City != nil {
		return u.City.ID
	}
	return ""
}

func (u UserState) GetServiceID() string {
	if u.Service != nil {
		return u.Service.ID
	}
	return ""
}

type City struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ServiceCategory struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	ID         string `json:"id"`
	Name       string `json:"name" validate:"required"`
	CategoryID string `json:"category_id" validate:"required"`
}

type Master struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Images      []string `json:"images"`
	Description string   `json:"description"`
	Contact     string   `json:"contact"`
	CityID      string   `json:"city_id"`
}

type MasterRegForm struct {
	ID                string   `json:"id,omitempty"`
	Name              string   `json:"name" validate:"required"`
	Images            []string `json:"images" validate:"required"`
	Description       string   `json:"description,omitempty"`
	Contact           string   `json:"contact" validate:"required"`
	CityID            string   `json:"city_id" validate:"required"`
	ServiceCategoryID string   `json:"service_category_id" validate:"required"`
	ServiceIDs        []string `json:"service_ids" validate:"required"`
}
