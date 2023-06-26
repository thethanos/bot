package bot

import (
	"fmt"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/entities"
	"multimessenger_bot/internal/logger"
	ma "multimessenger_bot/internal/messenger_adapter"
	"strings"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

type StepBase struct {
	logger     logger.Logger
	inProgress bool
	state      *entities.UserState
	dbAdapter  *db_adapter.DbAdapter
}

func (s *StepBase) IsInProgress() bool {
	return s.inProgress
}

func (s *StepBase) Reset() {
}

func (s *StepBase) SetInProgress(flag bool) {
	s.inProgress = flag
}

type YesNo struct {
	StepBase
	question Question
	yesStep  StepType
	noStep   StepType
}

func (y *YesNo) Request(msg *ma.Message) *ma.Message {
	y.logger.Infof("YesNo step is sending request")
	y.inProgress = true
	if msg.Source == ma.TELEGRAM {
		keyboard := makeKeyboard([]string{"Да", "Нет"})
		return ma.NewTextMessage(y.question.Text, msg, keyboard, false)
	}
	return ma.NewTextMessage("this messenger is unsupported yet", msg, nil, true)
}

func (y *YesNo) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	y.logger.Infof("YesNo step is processing response")
	y.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "да" || userAnswer == "1" {
		y.logger.Infof("Next step is %s", getStepTypeName(y.yesStep))
		return nil, y.yesStep
	}
	y.logger.Infof("Next step is %s", getStepTypeName(y.yesStep))
	return nil, y.noStep
}

type Prompt struct {
	StepBase
	question Question
	nextStep StepType
}

func (p *Prompt) Request(msg *ma.Message) *ma.Message {
	p.logger.Infof("Prompt step is sending request")
	p.inProgress = true
	if msg.Source == ma.TELEGRAM {
		keyboard := makeKeyboard([]string{"Назад"})
		return ma.NewTextMessage(p.question.Text, msg, keyboard, false)
	}
	return ma.NewTextMessage("this messenger is unsupported yet", msg, nil, true)
}

func (p *Prompt) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	p.logger.Infof("Prompt step is processing response")
	p.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		p.logger.Info("Next step is PreviousStep")
		return nil, PreviousStep
	}
	p.state.RawInput[p.question.Field] = msg.Text
	p.logger.Infof("Next step is %s", getStepTypeName(p.nextStep))
	return nil, p.nextStep
}

type MainMenu struct {
	StepBase
}

func (m *MainMenu) Request(msg *ma.Message) *ma.Message {
	m.logger.Infof("MainMenu step is sending request")
	m.state.Reset()
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Город"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Услуги"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Поиск моделей"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "По вопросам сотрудничества"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		m.inProgress = true
		return ma.NewTextMessage("Главное меню", msg, keyboard, false)
	}
	return ma.NewTextMessage("this messenger is unsupported yet", msg, nil, true)
}

func (m *MainMenu) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	m.logger.Infof("MainMenu step is processing response")
	m.inProgress = false

	switch strings.ToLower(msg.Text) {
	case "город":
		m.logger.Infof("Next step is MainMenuCitySelectionStep")
		return nil, MainMenuCitySelectionStep
	case "услуги":
		m.logger.Infof("Next step is MainMenuServiceCategorySelectionStep")
		return nil, MainMenuServiceCategorySelectionStep
	case "поиск моделей":
		m.logger.Infof("Next step is FindModelStep")
		return nil, FindModelStep
	case "по вопросам сотрудничества":
		m.logger.Infof("Next step is CollaborationStep")
		return nil, CollaborationStep
	case "админ":
		m.logger.Infof("Next step is AdminStep")
		return nil, AdminStep
	}

	return ma.NewTextMessage("Пожалуйста выберите ответ из списка.", msg, nil, false), EmptyStep
}

type MasterSelection struct {
	StepBase
}

func (m *MasterSelection) Request(msg *ma.Message) *ma.Message {
	m.logger.Info("MasterSelection step is sending request")
	m.inProgress = true
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Каталог мастеров", WebApp: &tgbotapi.WebAppInfo{
			Url: fmt.Sprintf("https://bot-dev-domain.com/pages/index.html?city_id=%s&service_id=%s", m.state.GetCityID(), m.state.GetServiceID()),
		}}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться назад"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться на главную"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage("Выбор мастера", msg, keyboard, false)
	}
	return ma.NewTextMessage("this messenger is unsupported yet", msg, nil, true)
}

func (m *MasterSelection) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	m.logger.Infof("MasterSelection step is processing response")
	m.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "вернуться назад" {
		m.logger.Infof("Next step is PreviousStep")
		return nil, PreviousStep
	}
	if userAnswer == "вернуться на главную" {
		m.logger.Infof("Next step is MainMenuStep")
		return nil, MainMenuStep
	}

	return nil, EmptyStep
}

type FindModel struct {
	StepBase
}

func (f *FindModel) Request(msg *ma.Message) *ma.Message {
	f.logger.Info("FindModel step is sending request")
	f.inProgress = true
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Каталог моделей", WebApp: &tgbotapi.WebAppInfo{
			Url: fmt.Sprintf("https://bot-dev-domain.com/masters?city_id=%s&service_id=%s", f.state.GetCityID(), f.state.GetServiceID()),
		}}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться на главную"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage("Поиск моделей", msg, keyboard, false)
	}
	return nil
}

func (f *FindModel) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	f.logger.Infof("FindModel step is processing response")
	f.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "вернуться на главную" {
		f.logger.Infof("Next step is MainMenuStep")
		return nil, MainMenuStep
	}
	return nil, EmptyStep
}

type Collaboration struct {
	StepBase
}

func (c *Collaboration) Request(msg *ma.Message) *ma.Message {
	c.logger.Info("Collaboration step is sending request")
	c.inProgress = true
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться на главную"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage("Всем привет! Меня зовут Маша и я алкоголик. Давайте сотрудничать.", msg, keyboard, false)
	}
	return nil
}

func (c *Collaboration) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	c.logger.Infof("Collaboration step is processing response")
	c.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "вернуться на главную" {
		c.logger.Infof("Next step is MainMenuStep")
		return nil, PreviousStep
	}
	return nil, EmptyStep
}

type Admin struct {
	StepBase
}

func (a *Admin) Request(msg *ma.Message) *ma.Message {
	a.logger.Info("Admin step is sending request")
	a.inProgress = true
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Добавить категорию услуг"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Добавить услугу"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Добавить город"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Добавить мастера", WebApp: &tgbotapi.WebAppInfo{
			Url: "https://bot-dev-domain.com/pages/registration_form.html"},
		}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться на главную"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage("Панель управления", msg, keyboard, false)
	}
	return ma.NewTextMessage("this messenger is unsupported yet", msg, nil, true)
}

func (a *Admin) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	a.logger.Infof("Admin step is processing response")
	a.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "вернуться на главную" {
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
		return ma.NewTextMessage("Пожалуйста выберите ответ из списка.", msg, nil, false), EmptyStep
	}
}
