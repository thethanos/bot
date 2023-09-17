package bot

import (
	"bot/internal/dbadapter"
	"bot/internal/entities"
	ma "bot/internal/msgadapter"
	"fmt"
	"strings"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

type ServiceCategoryStepMode interface {
	GetServCategories(cityId uint) ([]*entities.ServiceCategory, error)
	Text() string
	Buttons() [][]tgbotapi.KeyboardButton
	NextStep() StepType
}

type BaseServiceCategoryMode struct {
	dbAdapter *dbadapter.DBAdapter
}

func (b *BaseServiceCategoryMode) GetServCategories(cityId uint) ([]*entities.ServiceCategory, error) {
	return b.dbAdapter.GetServCategories(cityId, 0, -1)
}

func (b *BaseServiceCategoryMode) Text() string {
	return "Выберите категорию"
}

func (b *BaseServiceCategoryMode) Buttons() [][]tgbotapi.KeyboardButton {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться назад"}})
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться на главную"}})
	return rows
}

func (b *BaseServiceCategoryMode) NextStep() StepType {
	return ServiceSelectionStep
}

type MainMenuServiceCategoryMode struct {
	BaseServiceCategoryMode
}

func (m *MainMenuServiceCategoryMode) GetServCategories(cityId uint) ([]*entities.ServiceCategory, error) {
	return m.dbAdapter.GetServCategories(0, 0, -1)
}

func (m *MainMenuServiceCategoryMode) Text() string {
	return "По услуге"
}

func (m *MainMenuServiceCategoryMode) Buttons() [][]tgbotapi.KeyboardButton {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться на главную"}})
	return rows
}

func (m *MainMenuServiceCategoryMode) NextStep() StepType {
	return MainMenuServiceSelectionStep
}

type ServiceCategorySelection struct {
	StepBase
	categories []*entities.ServiceCategory
	mode       ServiceCategoryStepMode
}

func (s *ServiceCategorySelection) Request(msg *ma.Message) *ma.Message {
	s.logger.Infof("ServiceCategorySelection step is sending request")
	s.inProgress = true

	categories, _ := s.mode.GetServCategories(s.state.GetCityID())

	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		for _, category := range categories {
			rows = append(rows, []tgbotapi.KeyboardButton{{Text: category.Name}})
		}
		rows = append(rows, s.mode.Buttons()...)
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(categories) == 0 {
			return ma.NewTextMessage("По вашему запросу ничего не найдено", msg, keyboard, false)
		}

		s.categories = categories
		return ma.NewTextMessage(s.mode.Text(), msg, keyboard, false)
	}
	return ma.NewTextMessage("this messenger is unsupported yet", msg, nil, true)
}

func (s *ServiceCategorySelection) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	s.logger.Info("ServiceCategorySelection step is processing response")
	s.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "вернуться назад" {
		s.logger.Info("Next step is PreviousStep")
		return nil, PreviousStep
	}
	if userAnswer == "вернуться на главную" {
		s.logger.Info("Next step is MainMenuStep")
		return nil, MainMenuStep
	}

	for idx, service := range s.categories {
		if userAnswer == strings.ToLower(service.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			s.state.ServiceCategory = service
			s.logger.Infof("Next step is %s", getStepTypeName(s.mode.NextStep()))
			return nil, s.mode.NextStep()
		}
	}

	s.inProgress = true
	s.logger.Info("Next step is EmptyStep")
	return ma.NewTextMessage("Пожалуйста выберите ответ из списка.", msg, nil, false), EmptyStep
}

func (s *ServiceCategorySelection) Reset() {
	s.state.ServiceCategory = nil
}

type ServiceSelectionStepMode interface {
	GetServicesList(categoryId, cityId uint) ([]*entities.Service, error)
	MenuItems([]*entities.Service) [][]tgbotapi.KeyboardButton
	Buttons() [][]tgbotapi.KeyboardButton
	NextStep() StepType
}

type BaseServiceSelectionMode struct {
	dbAdapter *dbadapter.DBAdapter
}

func (b *BaseServiceSelectionMode) GetServicesList(categoryId, cityId uint) ([]*entities.Service, error) {
	return b.dbAdapter.GetServices(categoryId, cityId, 0, -1)
}

func (b *BaseServiceSelectionMode) MenuItems(services []*entities.Service) [][]tgbotapi.KeyboardButton {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	for _, service := range services {
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: service.Name}})
	}
	return rows
}

func (b *BaseServiceSelectionMode) Buttons() [][]tgbotapi.KeyboardButton {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться назад"}})
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться на главную"}})
	return rows
}

func (b *BaseServiceSelectionMode) NextStep() StepType {
	return MasterSelectionStep
}

type MainMenuServiceSelectionMode struct {
	BaseServiceSelectionMode
}

func (m *MainMenuServiceSelectionMode) GetServicesList(categoryId, cityId uint) ([]*entities.Service, error) {
	return m.dbAdapter.GetServices(categoryId, 0, 0, -1)
}

func (m *MainMenuServiceSelectionMode) NextStep() StepType {
	return CitySelectionStep
}

type ServiceSelection struct {
	StepBase
	services []*entities.Service
	mode     ServiceSelectionStepMode
}

func (s *ServiceSelection) Request(msg *ma.Message) *ma.Message {
	s.logger.Infof("ServiceSelection step is sending request")
	s.inProgress = true
	services, _ := s.mode.GetServicesList(s.state.ServiceCategory.ID, s.state.GetCityID())

	if msg.Source == ma.TELEGRAM {
		rows := s.mode.MenuItems(services)
		rows = append(rows, s.mode.Buttons()...)
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(services) == 0 {
			return ma.NewTextMessage("По вашему запросу ничего не найдено", msg, keyboard, false)
		}

		s.services = services
		return ma.NewTextMessage(s.state.ServiceCategory.Name, msg, keyboard, false)
	}
	return ma.NewTextMessage("this messenger is unsupported yet", msg, nil, true)
}

func (s *ServiceSelection) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	s.logger.Info("ServiceSelection step is processing response")
	s.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "вернуться назад" {
		s.logger.Info("Next step is PreviousStep")
		return nil, PreviousStep
	}
	if userAnswer == "вернуться на главную" {
		s.logger.Info("Next step is MainMenuStep")
		return nil, MainMenuStep
	}
	for idx, service := range s.services {
		if userAnswer == strings.ToLower(service.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			s.state.Service = service
			s.logger.Infof("Next step is %s", getStepTypeName(s.mode.NextStep()))
			return nil, s.mode.NextStep()
		}
	}

	s.inProgress = true
	s.logger.Infof("Next step is EmptyStep")
	return ma.NewTextMessage("Пожалуйста выберите ответ из списка.", msg, nil, false), EmptyStep
}

func (s *ServiceSelection) Reset() {
	s.state.Service = nil
}
