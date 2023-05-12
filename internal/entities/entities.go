package entities

type UserState struct {
	City     *City
	Service  *Service
	Master   *Master
	Cursor   int
	RawInput map[string]string
}

func (u *UserState) Reset() {
	u.City = nil
	u.Service = nil
	u.Master = nil
	u.Cursor = 0
	u.RawInput = make(map[string]string)
}

type City struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Service struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Master struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	CityID string `json:"city_id"`
}
