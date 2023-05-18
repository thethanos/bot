package bot

import (
	"fmt"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/entities"
	ma "multimessenger_bot/internal/messenger_adapter"
	"strings"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

type ServiceCategoryStepMode interface {
	GetServiceCategories(cityId string) ([]*entities.ServiceCategory, error)
	Text() string
	Buttons() [][]tgbotapi.KeyboardButton
	NextStep() StepType
}

type BaseServiceCategoryMode struct {
	dbAdapter *db_adapter.DbAdapter
}

func (b *BaseServiceCategoryMode) GetServiceCategories(cityId string) ([]*entities.ServiceCategory, error) {
	return b.dbAdapter.GetCategories(cityId)
}

func (b *BaseServiceCategoryMode) Text() string {
	return "Выберите услугу"
}

func (b *BaseServiceCategoryMode) Buttons() [][]tgbotapi.KeyboardButton {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
	return rows
}

func (b *BaseServiceCategoryMode) NextStep() StepType {
	return ServiceSelectionStep
}

type MainMenuServiceCategoryMode struct {
	BaseServiceCategoryMode
}

func (m *MainMenuServiceCategoryMode) GetServiceCategories(cityId string) ([]*entities.ServiceCategory, error) {
	return m.dbAdapter.GetCategories("")
}

func (m *MainMenuServiceCategoryMode) Text() string {
	return "По услуге"
}

func (m *MainMenuServiceCategoryMode) Buttons() [][]tgbotapi.KeyboardButton {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
	return rows
}

func (m *MainMenuServiceCategoryMode) NextStep() StepType {
	return MainMenuSericeSelectionStep
}

type MasterServiceCategoryMode struct {
	BaseServiceCategoryMode
}

func (m *MasterServiceCategoryMode) GetServiceCategories(cityId string) ([]*entities.ServiceCategory, error) {
	return m.dbAdapter.GetCategories("")
}

type AdminServiceCategoryMode struct {
	BaseServiceCategoryMode
}

func (a *AdminServiceCategoryMode) GetServiceCategories(cityId string) ([]*entities.ServiceCategory, error) {
	return a.dbAdapter.GetCategories("")
}

func (a *AdminServiceCategoryMode) NextStep() StepType {
	return AddServiceStep
}

type ServiceCategorySelection struct {
	StepBase
	categories []*entities.ServiceCategory
	mode       ServiceCategoryStepMode
}

func (s *ServiceCategorySelection) Request(msg *ma.Message) *ma.Message {
	s.logger.Infof("ServiceCategorySelection step is sending request")
	s.inProgress = true

	categories, _ := s.mode.GetServiceCategories(s.state.GetCityID())

	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		for _, category := range categories {
			rows = append(rows, []tgbotapi.KeyboardButton{{Text: category.Name}})
		}
		rows = append(rows, s.mode.Buttons()...)
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(categories) == 0 {
			return ma.NewMessage("Услуги не найдены", ma.REGULAR, msg, keyboard, nil)
		}

		s.categories = categories
		return ma.NewMessage(s.mode.Text(), ma.REGULAR, msg, keyboard, nil)
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
			s.state.ServiceCategory = service
			s.logger.Infof("Next step is %s", getStepTypeName(s.mode.NextStep()))
			return nil, s.mode.NextStep()
		}
	}

	s.inProgress = true
	s.logger.Info("Next step is EmptyStep")
	return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), EmptyStep
}

func (s *ServiceCategorySelection) Reset() {
	s.state.ServiceCategory = nil
}

type ServiceSelectionStepMode interface {
	NextStep() StepType
	Buttons() [][]tgbotapi.KeyboardButton
}

type BaseServiceSelectionMode struct {
}

func (b *BaseServiceSelectionMode) NextStep() StepType {
	return MasterSelectionStep
}

func (b *BaseServiceSelectionMode) Buttons() [][]tgbotapi.KeyboardButton {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
	return rows
}

type MainMenuServiceSelectionMode struct {
	BaseServiceSelectionMode
}

func (m *MainMenuServiceSelectionMode) NextStep() StepType {
	return CitySelectionStep
}

type RegistrationServiceSelectionMode struct {
	BaseServiceSelectionMode
}

func (b *RegistrationServiceSelectionMode) NextStep() StepType {
	return MasterRegistrationFinalStep
}

type ServiceSelection struct {
	StepBase
	services []*entities.Service
	mode     ServiceSelectionStepMode
}

func (s *ServiceSelection) Request(msg *ma.Message) *ma.Message {
	s.logger.Infof("ServiceSelection step is sending request")
	s.inProgress = true
	services, _ := s.dbAdapter.GetServices(s.state.ServiceCategory.ID)

	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		for _, service := range services {
			rows = append(rows, []tgbotapi.KeyboardButton{{Text: service.Name}})
		}
		rows = append(rows, s.mode.Buttons()...)
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(services) == 0 {
			return ma.NewMessage("По вашему запросу ничего не найдено", ma.REGULAR, msg, keyboard, nil)
		}

		s.services = services
		return ma.NewMessage(s.state.ServiceCategory.Name, ma.REGULAR, msg, keyboard, nil)
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
			s.state.Service = service
			s.logger.Infof("Next step is %s", getStepTypeName(s.mode.NextStep()))
			return nil, s.mode.NextStep()
		}
	}

	s.inProgress = true
	s.logger.Infof("Next step is EmptyStep")
	return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), EmptyStep
}

func (s *ServiceSelection) Reset() {
	s.state.Service = nil
}
