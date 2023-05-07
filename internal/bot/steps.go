package bot

import (
	"fmt"
	"multimessenger_bot/internal/db_adapter"
	ma "multimessenger_bot/internal/messenger_adapter"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type UserSession struct {
	CurrentStep Step
	State       UserState
}

type UserState struct {
	city    *db_adapter.City
	service *db_adapter.Service
	master  *db_adapter.Master
	cursor  int
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
	RegistrationStep
	RegistrationFinalStep
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

	if msg.Type == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 5)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Услуги"}}
		rows[1] = []tgbotapi.KeyboardButton{{Text: "Город"}}
		rows[2] = []tgbotapi.KeyboardButton{{Text: "Вопросы"}}
		rows[3] = []tgbotapi.KeyboardButton{{Text: "О нас"}}
		rows[4] = []tgbotapi.KeyboardButton{{Text: "Мастер"}}

		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		m.inProgress = true
		return &ma.Message{Text: "Главное меню", UserData: msg.UserData, Type: msg.Type, TgMarkup: keyboard}
	}

	text := "1) услуги\n2) город\n3) вопросы\n4) о нас\n5)мастер"
	m.inProgress = true
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (m *MainMenu) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	m.inProgress = false

	switch strings.ToLower(msg.Text) {
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

	if msg.Type == ma.TELEGRAM {

		rows := make([][]tgbotapi.KeyboardButton, len(cities))
		for idx, city := range cities {
			rows[idx] = make([]tgbotapi.KeyboardButton, 0)
			rows[idx] = append(rows[idx], tgbotapi.KeyboardButton{Text: city.Name})
		}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		c.cities = cities
		c.inProgress = true
		return &ma.Message{Text: " Выберите город", UserData: msg.UserData, Type: msg.Type, TgMarkup: keyboard}
	}

	text := ""
	for idx, city := range cities {
		text += fmt.Sprintf("%d. %s\n", idx+1, city.Name)
	}

	c.cities = cities
	c.inProgress = true
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (c *CitySelection) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	c.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	for idx, city := range c.cities {
		if userAnswer == strings.ToLower(city.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
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

	if msg.Type == ma.TELEGRAM {

		rows := make([][]tgbotapi.KeyboardButton, len(services))
		for idx, service := range services {
			rows[idx] = make([]tgbotapi.KeyboardButton, 0)
			rows[idx] = append(rows[idx], tgbotapi.KeyboardButton{Text: service.Name})
		}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		c.services = services
		c.inProgress = true
		return &ma.Message{Text: " Выберите услугу", UserData: msg.UserData, Type: msg.Type, TgMarkup: keyboard}
	}

	text := ""
	for idx, service := range services {
		text += fmt.Sprintf("%d. %s\n", idx+1, service.Name)
	}

	c.services = services
	c.inProgress = true
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (s *ServiceSelection) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	s.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	for idx, service := range s.services {
		if userAnswer == strings.ToLower(service.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			s.State.service = service
			if s.State.city == nil {
				return nil, CitySelectionStep
			} else {
				return nil, MasterSelectionStep
			}
		}
	}

	s.inProgress = true
	return &ma.Message{Text: "Пожалуйста выберите ответ из списка.", UserData: msg.UserData, Type: msg.Type}, EmptyStep
}

func (s *ServiceSelection) IsInProgress() bool {
	return s.inProgress
}

type MasterSelection struct {
	StepBase
	masters []*db_adapter.Master
}

func (m *MasterSelection) Request(msg *ma.Message) *ma.Message {

	masters, _ := m.DbAdapter.GetMasters(m.State.city, m.State.service)

	if msg.Type == ma.TELEGRAM {

		rows := make([][]tgbotapi.KeyboardButton, len(masters))
		for idx, master := range masters {
			rows[idx] = make([]tgbotapi.KeyboardButton, 0)
			rows[idx] = append(rows[idx], tgbotapi.KeyboardButton{Text: master.Name})
		}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		m.masters = masters
		m.inProgress = true
		return &ma.Message{Text: " Выберите мастера", UserData: msg.UserData, Type: msg.Type, TgMarkup: keyboard}
	}

	text := ""
	for idx, master := range masters {
		text += fmt.Sprintf("%d. %s", idx+1, master.Name)
	}

	m.masters = masters
	m.inProgress = true
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (m *MasterSelection) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	m.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	for idx, master := range m.masters {
		if userAnswer == strings.ToLower(master.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
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

	if msg.Type == ma.TELEGRAM {

		rows := make([][]tgbotapi.KeyboardButton, 2)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Да"}}
		rows[1] = []tgbotapi.KeyboardButton{{Text: "Нет"}}

		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		text := fmt.Sprintf("Ваша запись\nУслуга: %s\nГород: %s\nМастер: %s\nПодтвердить?",
			f.State.service.Name,
			f.State.city.Name,
			f.State.master.Name,
		)

		f.inProgress = true
		return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type, TgMarkup: keyboard}
	}

	text := fmt.Sprintf("Ваша запись\nУслуга: %s\nГород: %s\nМастер: %s\nПодтвердить?\nДа\nНет",
		f.State.service.Name,
		f.State.city.Name,
		f.State.master.Name,
	)
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
