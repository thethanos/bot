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
	CatID   uint   `json:"catID" validate:"required"`
	CatName string `json:"catName"`
}
