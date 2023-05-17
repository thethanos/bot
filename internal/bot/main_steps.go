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
		rows := make([][]tgbotapi.KeyboardButton, 6)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Список услуг"}}
		rows[1] = []tgbotapi.KeyboardButton{{Text: "По городу"}}
		rows[2] = []tgbotapi.KeyboardButton{{Text: "О нас"}}
		rows[3] = []tgbotapi.KeyboardButton{{Text: "Для мастеров"}}
		rows[4] = []tgbotapi.KeyboardButton{{Text: "Модель"}}
		rows[5] = []tgbotapi.KeyboardButton{{Text: "Test"}}

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
		return nil, QuestionsStep
	case "для мастеров":
		return nil, AboutStep
	case "модель":
		return nil, MasterStep
	case "админ":
		return nil, AdminStep
	case "test":
		return nil, TestStep
	}

	return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), EmptyStep
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

	var cities []*entities.City
	if c.filter && c.State.Service != nil {
		cities, _ = c.DbAdapter.GetCities(c.State.Service.ID)
	} else {
		cities, _ = c.DbAdapter.GetCities("")
	}

	if msg.Source == ma.TELEGRAM {

		rows := make([][]tgbotapi.KeyboardButton, len(cities)+1)
		for idx, city := range cities {
			rows[idx] = make([]tgbotapi.KeyboardButton, 0)
			rows[idx] = append(rows[idx], tgbotapi.KeyboardButton{Text: city.Name})
		}
		rows[len(cities)] = make([]tgbotapi.KeyboardButton, 0)
		rows[len(cities)] = append(rows[len(cities)], tgbotapi.KeyboardButton{Text: "Назад"})
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
	return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), c.errStep
}

func (c *CitySelection) Reset() {
	c.State.City = nil
}

type CityPrompt struct {
	StepBase
}

func (c *CityPrompt) Request(msg *ma.Message) *ma.Message {
	c.logger.Infof("CityPrompt step is sending request")
	c.inProgress = true

	text := "Введите город"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 1)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Главное меню"}}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true, OneTimeKeyboard: true}
		return ma.NewMessage(text, ma.REGULAR, msg, keyboard, nil)
	}

	return ma.NewMessage(fmt.Sprintf("%s\n1. Назад\n2. Главное меню", text), ma.REGULAR, msg, nil, nil)
}

func (c *CityPrompt) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	c.logger.Infof("CityPrompt step is processing response")
	c.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "главное меню" || userAnswer == "1" {
		return nil, MainMenuStep
	}

	city, err := c.DbAdapter.GetCity(msg.Text)
	if err != nil {
		c.inProgress = true
		c.logger.Infof("Next step is EmptyStep")
		return ma.NewMessage(fmt.Sprintf("По запросу %s ничего не найдено", msg.Text), ma.REGULAR, msg, nil, nil), EmptyStep
	}
	c.State.City = city
	if c.State.ServiceCategory == nil {
		c.logger.Infof("Next step is ServiceSelectionStep")
		return nil, ServiceCategorySelectionStep
	}
	c.logger.Infof("Next step is MasterSelectionStep")
	return nil, MasterSelectionStep
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
		rows := make([][]tgbotapi.KeyboardButton, len(categories)+1)
		for idx, category := range categories {
			rows[idx] = make([]tgbotapi.KeyboardButton, 0)
			rows[idx] = append(rows[idx], tgbotapi.KeyboardButton{Text: category.Name})
		}
		rows[len(categories)] = make([]tgbotapi.KeyboardButton, 0)
		rows[len(categories)] = append(rows[len(categories)], tgbotapi.KeyboardButton{Text: "Назад"})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(categories) == 0 {
			return ma.NewMessage("Услуги не найдены", ma.REGULAR, msg, keyboard, nil)
		}

		s.categories = categories
		return ma.NewMessage(" Выберите категорию услуг", ma.REGULAR, msg, keyboard, nil)
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
}

func (s *ServiceSelection) Request(msg *ma.Message) *ma.Message {
	s.logger.Infof("ServiceSelection step is sending request")
	s.inProgress = true
	services, _ := s.DbAdapter.GetServices(s.State.ServiceCategory.ID)

	if msg.Source == ma.TELEGRAM {

		rows := make([][]tgbotapi.KeyboardButton, len(services)+1)
		for idx, service := range services {
			rows[idx] = make([]tgbotapi.KeyboardButton, 0)
			rows[idx] = append(rows[idx], tgbotapi.KeyboardButton{Text: service.Name})
		}
		rows[len(services)] = make([]tgbotapi.KeyboardButton, 0)
		rows[len(services)] = append(rows[len(services)], tgbotapi.KeyboardButton{Text: "Назад"})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(services) == 0 {
			return ma.NewMessage("По вашему запросу ничего не найдено", ma.REGULAR, msg, keyboard, nil)
		}

		s.services = services
		return ma.NewMessage(" Выберите услугу", ma.REGULAR, msg, keyboard, nil)
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
		s.logger.Infof("Next step is PreviousStep")
		return nil, PreviousStep
	}
	for idx, service := range s.services {
		if userAnswer == strings.ToLower(service.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			s.State.Service = service
			if s.State.City == nil {
				s.logger.Infof("Next step is CityPromptStep")
				return nil, CityPromptStep
			}
			s.logger.Infof("Next step is MasterSelectionStep")
			return nil, MasterSelectionStep
		}
	}

	s.inProgress = true
	s.logger.Infof("Next step is EmptyStep")
	return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), EmptyStep
}

type MasterSelection struct {
	StepBase
	masters []*entities.Master
}

func (m *MasterSelection) Request(msg *ma.Message) *ma.Message {
	m.logger.Infof("MasterSelection step is sending request")
	m.inProgress = true
	masters, _ := m.DbAdapter.GetMasters(m.State.City.ID, m.State.Service.ID)

	if msg.Source == ma.TELEGRAM {

		rows := make([][]tgbotapi.KeyboardButton, len(masters)+1)
		for idx, master := range masters {
			rows[idx] = make([]tgbotapi.KeyboardButton, 0)
			rows[idx] = append(rows[idx], tgbotapi.KeyboardButton{Text: master.Name})
		}
		rows[len(masters)] = make([]tgbotapi.KeyboardButton, 0)
		rows[len(masters)] = append(rows[len(masters)], tgbotapi.KeyboardButton{Text: "Назад"})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(masters) == 0 {
			return ma.NewMessage("По вашему запросу ничего не найдено", ma.REGULAR, msg, keyboard, nil)
		}

		m.masters = masters
		return ma.NewMessage(" Выберите мастера", ma.REGULAR, msg, keyboard, nil)
	}

	text := ""
	for idx, master := range masters {
		text += fmt.Sprintf("%d. %s", idx+1, master.Name)
	}

	m.masters = masters
	return ma.NewMessage(text, ma.REGULAR, msg, nil, nil)
}

func (m *MasterSelection) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	if msg.Type == ma.CALLBACK {
		return nil, EmptyStep
	}
	m.logger.Infof("MasterSelection step is processing response")
	m.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" || userAnswer == fmt.Sprintf("%d", len(m.masters)+1) {
		m.logger.Infof("Next step is PreviousStep")
		return nil, PreviousStep
	}
	for idx, master := range m.masters {
		if userAnswer == strings.ToLower(master.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			m.State.Master = master
			m.logger.Infof("Next step is FinalStep")
			return nil, FinalStep
		}
	}

	m.inProgress = true
	m.logger.Infof("Next step is EmptyStep")
	return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), EmptyStep
}

type Final struct {
	StepBase
}

func (f *Final) Request(msg *ma.Message) *ma.Message {
	f.logger.Infof("Final step is sending request")
	f.inProgress = true
	text := fmt.Sprintf("Ваша запись\nУслуга: %s\nГород: %s\nМастер: %s\n\nПодтвердить?",
		f.State.Service.Name,
		f.State.City.Name,
		f.State.Master.Name,
	)

	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 2)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Да"}}
		rows[1] = []tgbotapi.KeyboardButton{{Text: "Нет"}}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewMessage(text, ma.REGULAR, msg, keyboard, nil)
	}
	return ma.NewMessage(fmt.Sprintf("%s\n1. Да\n2. Нет", text), ma.REGULAR, msg, nil, nil)
}

func (f *Final) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	if msg.Type == ma.CALLBACK {
		return nil, EmptyStep
	}
	f.logger.Infof("Final step is processing response")
	f.inProgress = false

	switch msg.Text {
	case "Да":
		f.State.Reset()
		f.logger.Infof("Next step is MainMenuRequestStep")
		return ma.NewMessage("Запись завершена", ma.REGULAR, msg, nil, nil), MainMenuRequestStep
	case "Нет":
		f.State.Reset()
		f.logger.Infof("Next step is MainMenuRequestStep")
		return ma.NewMessage("Запись отменена", ma.REGULAR, msg, nil, nil), MainMenuRequestStep
	default:
		f.inProgress = true
		f.logger.Infof("Next step is EmptyStep")
		return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), EmptyStep
	}
}

type Admin struct {
	StepBase
}

func (a *Admin) Request(msg *ma.Message) *ma.Message {
	a.logger.Info("Admin step is sending request")
	a.inProgress = true

	text := "Панель управления"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 4)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Добавить категорию услуг"}}
		rows[1] = []tgbotapi.KeyboardButton{{Text: "Добавить услугу"}}
		rows[2] = []tgbotapi.KeyboardButton{{Text: "Добавить город"}}
		rows[3] = []tgbotapi.KeyboardButton{{Text: "Назад"}}
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
