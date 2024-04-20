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
	ID      string `json:"id"`
	Name    string `json:"name" validate:"required"`
	CatID   string `json:"catID" validate:"required"`
	CatName string `json:"catName"`
}
