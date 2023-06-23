package entities

type UserState struct {
	City            *City
	ServiceCategory *ServiceCategory
	Service         *Service
	Master          *Master
	Cursor          int
	RawInput        map[string]string
}

func (u *UserState) Reset() {
	u.City = nil
	u.ServiceCategory = nil
	u.Service = nil
	u.Master = nil
	u.Cursor = 0
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
	ID       string `json:"id"`
	IndexStr string `json:"-"`
	Name     string `json:"name"`
}

type ServiceCategory struct {
	ID       string `json:"id"`
	IndexStr string `json:"-"`
	Name     string `json:"name"`
}

type Service struct {
	ID         string `json:"id"`
	IndexStr   string `json:"-"`
	Name       string `json:"name"`
	CategoryID string `json:"-"`
}

type Master struct {
	ID          string   `json:"id"`
	IndexStr    string   `json:"index_str"`
	Name        string   `json:"name"`
	Images      []string `json:"images"`
	Description string   `json:"description"`
	Contact     string   `json:"contact"`
	CityID      string   `json:"city_id"`
}
