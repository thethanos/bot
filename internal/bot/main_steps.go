package bot

import (
	"fmt"
	"multimessenger_bot/internal/entities"
	ma "multimessenger_bot/internal/messenger_adapter"
	"strings"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

type MainMenu struct {
	StepBase
}

func (m *MainMenu) Request(msg *ma.Message) *ma.Message {
	m.logger.Infof("MainMenu step is sending request")
	m.State.Reset()
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 5)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Список услуг"}}
		rows[1] = []tgbotapi.KeyboardButton{{Text: "По городу"}}
		rows[2] = []tgbotapi.KeyboardButton{{Text: "О нас"}}
		rows[3] = []tgbotapi.KeyboardButton{{Text: "Для мастеров"}}
		rows[4] = []tgbotapi.KeyboardButton{{Text: "Модель"}}

		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		m.inProgress = true
		return ma.NewMessage("Главное меню", ma.REGULAR, msg, keyboard, nil)
	}

	text := "1. услуги\n2. город\n3. вопросы\n4. о нас\n5. мастер"
	m.inProgress = true
	return ma.NewMessage(text, ma.REGULAR, msg, nil, nil)
}

func (m *MainMenu) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	m.logger.Infof("MainMenu step is processing response")
	if msg.Type == ma.CALLBACK {
		return nil, EmptyStep
	}
	m.inProgress = false

	switch strings.ToLower(msg.Text) {
	case "список услуг":
		return nil, ServiceCategorySelectionStep
	case "по городу":
		return nil, CityPromptStep
	case "о нас":
		return nil, AboutStep
	case "для мастеров":
		return nil, MasterStep
	case "модель":
		return nil, EmptyStep
	case "админ":
		return nil, AdminStep
	}

	return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), EmptyStep
}

type ServiceCategorySelection struct {
	StepBase
	categories []*entities.ServiceCategory
	filter     bool
	addSrvMode bool
	errStep    StepType
}

func (s *ServiceCategorySelection) Request(msg *ma.Message) *ma.Message {
	s.logger.Infof("ServiceCategorySelection step is sending request")
	s.inProgress = true

	var categories []*entities.ServiceCategory
	if s.filter && s.State.City != nil {
		categories, _ = s.DbAdapter.GetCategories(s.State.City.ID)
	} else {
		categories, _ = s.DbAdapter.GetCategories("")
	}

	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		for _, category := range categories {
			rows = append(rows, []tgbotapi.KeyboardButton{{Text: category.Name}})
		}
		if s.State.City != nil {
			rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		}
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(categories) == 0 {
			return ma.NewMessage("Услуги не найдены", ma.REGULAR, msg, keyboard, nil)
		}

		s.categories = categories
		return ma.NewMessage("По услуге", ma.REGULAR, msg, keyboard, nil)
	}

	text := ""
	for idx, category := range categories {
		text += fmt.Sprintf("%d. %s\n", idx+1, category.Name)
	}
	text += fmt.Sprintf("%d. Назад\n", len(categories)+1)

	s.categories = categories
	return ma.NewMessage(text, ma.REGULAR, msg, nil, nil)
}

func (s *ServiceCategorySelection) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	s.logger.Info("ServiceCategorySelection step is processing response")
	if msg.Type == ma.CALLBACK {
		return nil, EmptyStep
	}
	s.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" || userAnswer == fmt.Sprintf("%d", len(s.categories)+1) {
		return nil, PreviousStep
	}
	if userAnswer == "главное меню" {
		return nil, MainMenuStep
	}

	for idx, service := range s.categories {
		if userAnswer == strings.ToLower(service.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			s.State.ServiceCategory = service
			if s.addSrvMode {
				s.logger.Info("Next step is AddServiceStep")
				return nil, AddServiceStep
			}
			s.logger.Info("Next step is ServiceSelectionStep")
			return nil, ServiceSelectionStep
		}
	}

	s.inProgress = true
	s.logger.Infof("Next step is %s", getStepTypeName(s.errStep))
	return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), s.errStep
}

func (s *ServiceCategorySelection) Reset() {
	s.State.ServiceCategory = nil
}

type ServiceSelection struct {
	StepBase
	services []*entities.Service
	nextStep StepType
}

func (s *ServiceSelection) Request(msg *ma.Message) *ma.Message {
	s.logger.Infof("ServiceSelection step is sending request")
	s.inProgress = true
	services, _ := s.DbAdapter.GetServices(s.State.ServiceCategory.ID)

	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		for _, service := range services {
			rows = append(rows, []tgbotapi.KeyboardButton{{Text: service.Name}})
		}
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(services) == 0 {
			return ma.NewMessage("По вашему запросу ничего не найдено", ma.REGULAR, msg, keyboard, nil)
		}

		s.services = services
		return ma.NewMessage(s.State.ServiceCategory.Name, ma.REGULAR, msg, keyboard, nil)
	}

	text := ""
	for idx, service := range services {
		text += fmt.Sprintf("%d. %s", idx+1, service.Name)
	}

	s.services = services
	return ma.NewMessage(text, ma.REGULAR, msg, nil, nil)
}

func (s *ServiceSelection) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	s.logger.Info("ServiceSelection step is processing response")
	s.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" || userAnswer == fmt.Sprintf("%d", len(s.services)+1) {
		return nil, PreviousStep
	}
	if userAnswer == "главное меню" {
		return nil, MainMenuStep
	}
	for idx, service := range s.services {
		if userAnswer == strings.ToLower(service.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			s.State.Service = service
			if s.State.City == nil {
				s.logger.Infof("Next step is CitySelectionStep")
				return nil, CitySelectionStep
			}
			s.logger.Infof("Next step is MasterSelectionStep")
			return nil, s.nextStep
		}
	}

	s.inProgress = true
	s.logger.Infof("Next step is EmptyStep")
	return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), EmptyStep
}

func (s *ServiceSelection) Reset() {
	s.State.Service = nil
}

type CityPrompt struct {
	StepBase
	nextStep StepType
}

func (c *CityPrompt) Request(msg *ma.Message) *ma.Message {
	c.logger.Infof("CityPrompt step is sending request")
	c.inProgress = true

	text := "Введите город"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		if c.State.ServiceCategory != nil {
			rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		}
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true, OneTimeKeyboard: true}
		return ma.NewMessage(text, ma.REGULAR, msg, keyboard, nil)
	}

	return ma.NewMessage(fmt.Sprintf("%s\n1. Назад\n2. Главное меню", text), ma.REGULAR, msg, nil, nil)
}

func (c *CityPrompt) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	c.logger.Infof("CityPrompt step is processing response")
	c.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		return nil, PreviousStep
	}
	if userAnswer == "главное меню" {
		return nil, MainMenuStep
	}

	city, err := c.DbAdapter.GetCity(msg.Text)
	if err != nil {
		c.inProgress = true
		c.logger.Infof("Next step is CityPromptStep")
		return ma.NewMessage(fmt.Sprintf("По запросу %s ничего не найдено", msg.Text), ma.REGULAR, msg, nil, nil), CityPromptStep
	}

	c.State.City = city
	if c.State.ServiceCategory == nil {
		c.logger.Infof("Next step is %s", getStepTypeName(c.nextStep))
		return nil, c.nextStep
	}
	c.logger.Infof("Next step is MasterSelectionStep")
	return nil, EmptyStep
}

func (c *CityPrompt) Reset() {
	c.State.City = nil
}

type CitySelection struct {
	StepBase
	cities       []*entities.City
	filter       bool
	checkService bool
	nextStep     StepType
	errStep      StepType
}

func (c *CitySelection) Request(msg *ma.Message) *ma.Message {
	c.logger.Infof("CitySelection step is sending request")
	c.inProgress = true

	cities, _ := c.DbAdapter.GetCities(c.State.Service.ID)

	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		for _, city := range cities {
			rows = append(rows, []tgbotapi.KeyboardButton{{Text: city.Name}})
		}
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(cities) == 0 {
			return ma.NewMessage("По вашему запросу ничего не найдено", ma.REGULAR, msg, keyboard, nil)
		}

		c.cities = cities
		return ma.NewMessage(" Выберите город", ma.REGULAR, msg, keyboard, nil)
	}

	text := ""
	for idx, city := range cities {
		text += fmt.Sprintf("%d. %s\n", idx+1, city.Name)
	}
	text += fmt.Sprintf("%d. Назад\n", len(cities)+1)

	c.cities = cities
	return ma.NewMessage(text, ma.REGULAR, msg, nil, nil)
}

func (c *CitySelection) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	c.logger.Infof("CitySelection step is processing response")
	if msg.Type == ma.CALLBACK {
		return nil, EmptyStep
	}
	c.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		return nil, PreviousStep
	}
	if userAnswer == "главное меню" {
		return nil, MainMenuStep
	}

	for idx, city := range c.cities {
		if userAnswer == strings.ToLower(city.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			c.State.City = city
			return nil, EmptyStep
		}
	}

	c.inProgress = true
	return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), c.errStep
}

func (c *CitySelection) Reset() {
	c.State.City = nil
}

type Admin struct {
	StepBase
}

func (a *Admin) Request(msg *ma.Message) *ma.Message {
	a.logger.Info("Admin step is sending request")
	a.inProgress = true

	text := "Панель управления"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Добавить категорию услуг"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Добавить услугу"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Добавить город"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewMessage(text, ma.REGULAR, msg, keyboard, nil)
	}
	return ma.NewMessage(fmt.Sprintf("%s\n1. Добавить категорию услуг\n2: Добавить услугу\n3. Добавить город\n4. Назад", text), ma.REGULAR, msg, nil, nil)
}

func (a *Admin) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	if msg.Type == ma.CALLBACK {
		return nil, EmptyStep
	}
	a.logger.Infof("Admin step is processing response")
	a.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" || userAnswer == "3" {
		a.logger.Infof("Next step is PreviousStep")
		return nil, PreviousStep
	}

	switch userAnswer {
	case "добавить категорию услуг":
		a.logger.Info("Next step is AddServiceCategory")
		return nil, AddServiceCategoryStep
	case "добавить услугу":
		a.logger.Info("Next step is AddServiceStep")
		return nil, AdminServiceCategorySelectionStep
	case "добавить город":
		a.logger.Info("Next step is AddCityStep")
		return nil, AddCityStep
	default:
		a.inProgress = true
		a.logger.Info("Next step is EmptyStep")
		return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), EmptyStep
	}
}
