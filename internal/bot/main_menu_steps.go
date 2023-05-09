package bot

import (
	"fmt"
	"multimessenger_bot/internal/entities"
	ma "multimessenger_bot/internal/messenger_adapter"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MainMenu struct {
	StepBase
}

func (m *MainMenu) Request(msg *ma.Message) *ma.Message {
	m.State.Reset()
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

	text := "1. услуги\n2. город\n3. вопросы\n4. о нас\n5. мастер"
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

type CitySelection struct {
	StepBase
	cities       []*entities.City
	filter       bool
	checkService bool
	nextStep     int
	errStep      int
}

func (c *CitySelection) Request(msg *ma.Message) *ma.Message {
	c.inProgress = true

	var cities []*entities.City
	if c.filter && c.State.Service != nil {
		cities, _ = c.DbAdapter.GetCities(c.State.Service.ID)
	} else {
		cities, _ = c.DbAdapter.GetCities("")
	}

	if msg.Type == ma.TELEGRAM {

		rows := make([][]tgbotapi.KeyboardButton, len(cities)+1)
		for idx, city := range cities {
			rows[idx] = make([]tgbotapi.KeyboardButton, 0)
			rows[idx] = append(rows[idx], tgbotapi.KeyboardButton{Text: city.Name})
		}
		rows[len(cities)] = make([]tgbotapi.KeyboardButton, 0)
		rows[len(cities)] = append(rows[len(cities)], tgbotapi.KeyboardButton{Text: "Назад"})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(cities) == 0 {
			return &ma.Message{Text: "По вашему запросу ничего не найдено", UserData: msg.UserData, Type: msg.Type, TgMarkup: keyboard}
		}

		c.cities = cities
		return &ma.Message{Text: " Выберите город", UserData: msg.UserData, Type: msg.Type, TgMarkup: keyboard}
	}

	text := ""
	for idx, city := range cities {
		text += fmt.Sprintf("%d. %s\n", idx+1, city.Name)
	}
	text += fmt.Sprintf("%d. Назад\n", len(cities)+1)

	c.cities = cities
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (c *CitySelection) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	c.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" || userAnswer == fmt.Sprintf("%d", len(c.cities)+1) {
		return nil, PreviousStep
	}

	for idx, city := range c.cities {
		if userAnswer == strings.ToLower(city.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			c.State.City = city
			if c.checkService {
				if c.State.Service == nil {
					return nil, ServiceSelectionStep
				} else {
					return nil, c.nextStep
				}
			} else {
				return nil, c.nextStep
			}
		}
	}

	c.inProgress = true
	return &ma.Message{Text: "Пожалуйста выберите ответ из списка.", UserData: msg.UserData, Type: msg.Type}, c.errStep
}

func (c *CitySelection) Reset() {
	c.State.City = nil
}

type ServiceSelection struct {
	StepBase
	services  []*entities.Service
	filter    bool
	checkCity bool
	nextStep  int
	errStep   int
}

func (c *ServiceSelection) Request(msg *ma.Message) *ma.Message {
	c.inProgress = true

	var services []*entities.Service
	if c.filter && c.State.City != nil {
		services, _ = c.DbAdapter.GetServices(c.State.City.ID)
	} else {
		services, _ = c.DbAdapter.GetServices("")
	}

	if msg.Type == ma.TELEGRAM {

		rows := make([][]tgbotapi.KeyboardButton, len(services)+1)
		for idx, service := range services {
			rows[idx] = make([]tgbotapi.KeyboardButton, 0)
			rows[idx] = append(rows[idx], tgbotapi.KeyboardButton{Text: service.Name})
		}
		rows[len(services)] = make([]tgbotapi.KeyboardButton, 0)
		rows[len(services)] = append(rows[len(services)], tgbotapi.KeyboardButton{Text: "Назад"})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(services) == 0 {
			return &ma.Message{Text: "По вашему запросу ничего не найдено", UserData: msg.UserData, Type: msg.Type, TgMarkup: keyboard}
		}

		c.services = services
		return &ma.Message{Text: " Выберите услугу", UserData: msg.UserData, Type: msg.Type, TgMarkup: keyboard}
	}

	text := ""
	for idx, service := range services {
		text += fmt.Sprintf("%d. %s\n", idx+1, service.Name)
	}
	text += fmt.Sprintf("%d. Назад\n", len(services)+1)

	c.services = services
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (s *ServiceSelection) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	s.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" || userAnswer == fmt.Sprintf("%d", len(s.services)+1) {
		return nil, PreviousStep
	}
	for idx, service := range s.services {
		if userAnswer == strings.ToLower(service.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			s.State.Service = service
			if s.checkCity {
				if s.State.City == nil {
					return nil, CitySelectionStep
				} else {
					return nil, s.nextStep
				}
			} else {
				return nil, s.nextStep
			}
		}
	}

	s.inProgress = true
	return &ma.Message{Text: "Пожалуйста выберите ответ из списка.", UserData: msg.UserData, Type: msg.Type}, s.errStep
}

func (s *ServiceSelection) Reset() {
	s.State.Service = nil
}

type MasterSelection struct {
	StepBase
	masters []*entities.Master
}

func (m *MasterSelection) Request(msg *ma.Message) *ma.Message {
	m.inProgress = true

	masters, _ := m.DbAdapter.GetMasters(m.State.City.ID, m.State.Service.ID)

	if msg.Type == ma.TELEGRAM {

		rows := make([][]tgbotapi.KeyboardButton, len(masters)+1)
		for idx, master := range masters {
			rows[idx] = make([]tgbotapi.KeyboardButton, 0)
			rows[idx] = append(rows[idx], tgbotapi.KeyboardButton{Text: master.Name})
		}
		rows[len(masters)] = make([]tgbotapi.KeyboardButton, 0)
		rows[len(masters)] = append(rows[len(masters)], tgbotapi.KeyboardButton{Text: "Назад"})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(masters) == 0 {
			return &ma.Message{Text: "По вашему запросу ничего не найдено", UserData: msg.UserData, Type: msg.Type, TgMarkup: keyboard}
		}

		m.masters = masters
		return &ma.Message{Text: " Выберите мастера", UserData: msg.UserData, Type: msg.Type, TgMarkup: keyboard}
	}

	text := ""
	for idx, master := range masters {
		text += fmt.Sprintf("%d. %s", idx+1, master.Name)
	}

	m.masters = masters
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (m *MasterSelection) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	m.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" || userAnswer == fmt.Sprintf("%d", len(m.masters)+1) {
		return nil, PreviousStep
	}
	for idx, master := range m.masters {
		if userAnswer == strings.ToLower(master.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			m.State.Master = master
			return nil, FinalStep
		}
	}

	m.inProgress = true
	return &ma.Message{Text: "Пожалуйста выберите ответ из списка.", UserData: msg.UserData, Type: msg.Type}, EmptyStep
}

type Final struct {
	StepBase
}

func (f *Final) Request(msg *ma.Message) *ma.Message {

	text := fmt.Sprintf("Ваша запись\nУслуга: %s\nГород: %s\nМастер: %s\n\nПодтвердить?",
		f.State.Service.Name,
		f.State.City.Name,
		f.State.Master.Name,
	)

	if msg.Type == ma.TELEGRAM {

		rows := make([][]tgbotapi.KeyboardButton, 2)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Да"}}
		rows[1] = []tgbotapi.KeyboardButton{{Text: "Нет"}}

		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		f.inProgress = true
		return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type, TgMarkup: keyboard}
	}

	f.inProgress = true
	return &ma.Message{Text: fmt.Sprintf("%s\n1. Да\n2. Нет", text), UserData: msg.UserData, Type: msg.Type}
}

func (f *Final) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	f.inProgress = false

	switch msg.Text {
	case "Да":
		f.State.Reset()
		return &ma.Message{Text: "Запись завершена", UserData: msg.UserData, Type: msg.Type}, MainMenuRequestStep
	case "Нет":
		f.State.Reset()
		return &ma.Message{Text: "Запись отменена", UserData: msg.UserData, Type: msg.Type}, MainMenuRequestStep
	default:
		f.inProgress = true
		return &ma.Message{Text: "Пожалуйста выберите ответ из списка.", UserData: msg.UserData, Type: msg.Type}, EmptyStep
	}
}
