package bot

import (
	"fmt"
	"multimessenger_bot/internal/db_adapter"
	ma "multimessenger_bot/internal/messenger_adapter"
)

type UserSession struct {
	CurrentStep Step
	State       UserState
}

type UserState struct {
	city    *db_adapter.City
	service *db_adapter.Service
	master  *db_adapter.Master
}

const (
	MainMenuStep = iota
	ServiceSelectionStep
	CitySelectionStep
	QuestionsStep
	AboutStep
	MasterSelectionStep
	MasterStep
	FinalStep
	EmptyStep
)

type Step interface {
	ProcessResponse(*ma.Message) (*ma.Message, int)
	Request(*ma.Message) *ma.Message
	IsInProgress() bool
}

type StepBase struct {
	inProgress bool
	State      *UserState
	DbAdapter  *db_adapter.DbAdapter
}

type MainMenu struct {
	StepBase
}

func (m *MainMenu) Request(msg *ma.Message) *ma.Message {
	text := "1) услуги\n2) город\n3) вопросы\n4) о нас\n5)мастер"
	m.inProgress = true
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (m *MainMenu) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	m.inProgress = false

	switch msg.Text {
	case "услуги":
		return nil, ServiceSelectionStep
	case "город":
		return nil, CitySelectionStep
	case "вопросы":
		return nil, QuestionsStep
	case "о нас":
		return nil, AboutStep
	case "мастер":
		return nil, MasterStep
	}

	return &ma.Message{Text: "Пожалуйста выберите ответ из списка.", UserData: msg.UserData, Type: msg.Type}, EmptyStep
}

func (m *MainMenu) IsInProgress() bool {
	return m.inProgress
}

type CitySelection struct {
	StepBase
	cities []*db_adapter.City
}

func (c *CitySelection) Request(msg *ma.Message) *ma.Message {

	cities, _ := c.DbAdapter.GetCities(c.State.service)

	text := ""
	for _, city := range cities {
		text += fmt.Sprintf("%d) %s\n", city.ID, city.Name)
	}

	c.cities = cities
	c.inProgress = true
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (c *CitySelection) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	c.inProgress = false

	for _, city := range c.cities {
		if msg.Text == city.Name {
			c.State.city = city
			if c.State.service == nil {
				return nil, ServiceSelectionStep
			} else {
				return nil, MasterSelectionStep
			}
		}
	}

	c.inProgress = true
	return &ma.Message{Text: "Пожалуйста выберите ответ из списка.", UserData: msg.UserData, Type: msg.Type}, EmptyStep
}

func (c *CitySelection) IsInProgress() bool {
	return c.inProgress
}

type ServiceSelection struct {
	StepBase
	services []*db_adapter.Service
}

func (c *ServiceSelection) Request(msg *ma.Message) *ma.Message {

	services, _ := c.DbAdapter.GetServices(c.State.city)

	text := ""
	for _, service := range services {
		text += fmt.Sprintf("%d) %s\n", service.ID, service.Name)
	}

	c.services = services
	c.inProgress = true
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (c *ServiceSelection) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	c.inProgress = false

	for _, service := range c.services {
		if msg.Text == service.Name {
			c.State.service = service
			if c.State.city == nil {
				return nil, CitySelectionStep
			} else {
				return nil, MasterSelectionStep
			}
		}
	}

	c.inProgress = true
	return &ma.Message{Text: "Пожалуйста выберите ответ из списка.", UserData: msg.UserData, Type: msg.Type}, EmptyStep
}

func (c *ServiceSelection) IsInProgress() bool {
	return c.inProgress
}

type MasterSelection struct {
	StepBase
	masters []*db_adapter.Master
}

func (m *MasterSelection) Request(msg *ma.Message) *ma.Message {

	masters, _ := m.DbAdapter.GetMasters(m.State.city, m.State.service)
	text := ""
	for _, master := range masters {
		text += fmt.Sprintf("%d) %s", master.ID, master.Name)
	}

	m.masters = masters
	m.inProgress = true
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (m *MasterSelection) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	m.inProgress = false

	for _, master := range m.masters {
		if msg.Text == master.Name {
			m.State.master = master
			return nil, FinalStep
		}
	}

	m.inProgress = true
	return &ma.Message{Text: "Пожалуйста выберите ответ из списка.", UserData: msg.UserData, Type: msg.Type}, EmptyStep
}

func (m *MasterSelection) IsInProgress() bool {
	return m.inProgress
}

type Final struct {
	StepBase
}

func (f *Final) Request(msg *ma.Message) *ma.Message {
	text := fmt.Sprintf("Ваша запись\nУслуга: %s\nГород: %s\nМастер: %s\nПодтвердить?\nДа\nНет", f.State.service.Name, f.State.city.Name, f.State.master.Name)
	f.inProgress = true
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (f *Final) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	f.inProgress = false

	switch msg.Text {
	case "Да":
		f.State = &UserState{}
		return &ma.Message{Text: "Запись завершена", UserData: msg.UserData, Type: msg.Type}, MainMenuStep
	case "Нет":
		f.State = &UserState{}
		return &ma.Message{Text: "Запись отменена", UserData: msg.UserData, Type: msg.Type}, MainMenuStep
	default:
		f.inProgress = true
		return &ma.Message{Text: "Пожалуйста выберите ответ из списка.", UserData: msg.UserData, Type: msg.Type}, EmptyStep
	}
}

func (f *Final) IsInProgress() bool {
	return f.inProgress
}
