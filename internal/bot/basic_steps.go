package bot

import (
	"bot/internal/dbadapter"
	"bot/internal/entities"
	"bot/internal/logger"
	ma "bot/internal/msgadapter"
	"fmt"
	"strings"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

type StepBase struct {
	logger    logger.Logger
	state     *entities.UserState
	DBAdapter dbadapter.DBInterface
}

func (s *StepBase) Reset() {
}

type YesNo struct {
	StepBase
	question Question
	yesStep  StepType
	noStep   StepType
}

func (y *YesNo) Request(msg *ma.Message) *ma.Message {
	y.logger.Infof("YesNo step is sending request")
	if msg.Source == ma.TELEGRAM {
		keyboard := makeKeyboard([]string{"Да", "Нет"})
		return ma.NewTextMessage(y.question.Text, msg, keyboard, false)
	}
	return ma.NewTextMessage("this messenger is unsupported yet", msg, nil, true)
}

func (y *YesNo) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	y.logger.Infof("YesNo step is processing response")
	if Compare(msg.Text, "да") || Compare(msg.Text, "1") {
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
	if msg.Source == ma.TELEGRAM {
		keyboard := makeKeyboard([]string{Back})
		return ma.NewTextMessage(p.question.Text, msg, keyboard, false)
	}
	return ma.NewTextMessage("this messenger is unsupported yet", msg, nil, true)
}

func (p *Prompt) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	p.logger.Infof("Prompt step is processing response")
	if Compare(msg.Text, Back) {
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
		keyboard := makeKeyboard([]string{"Город", "Услуги", "Поиск моделей", "По вопросам сотрудничества"})
		return ma.NewTextMessage("Главное меню", msg, keyboard, false)
	}
	return ma.NewTextMessage("this messenger is unsupported yet", msg, nil, true)
}

func (m *MainMenu) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	m.logger.Infof("MainMenu step is processing response")

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
	}

	return ma.NewTextMessage("Пожалуйста выберите ответ из списка.", msg, nil, false), EmptyStep
}

type MasterSelection struct {
	StepBase
	GalleryURL string
}

func (m *MasterSelection) Request(msg *ma.Message) *ma.Message {
	m.logger.Info("MasterSelection step is sending request")
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Каталог мастеров", WebApp: &tgbotapi.WebAppInfo{
			Url: fmt.Sprintf("%s?city_id=%s&service_id=%s", m.GalleryURL, m.state.GetCityID(), m.state.GetServiceID()),
		}}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: Back}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: BackToMain}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage("Выбор мастера", msg, keyboard, false)
	}
	return ma.NewTextMessage("this messenger is unsupported yet", msg, nil, true)
}

func (m *MasterSelection) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	m.logger.Infof("MasterSelection step is processing response")
	if Compare(msg.Text, Back) {
		m.logger.Infof("Next step is PreviousStep")
		return nil, PreviousStep
	}
	if Compare(msg.Text, BackToMain) {
		m.logger.Infof("Next step is MainMenuStep")
		return nil, MainMenuStep
	}

	return nil, EmptyStep
}

type FindModel struct {
	StepBase
	ModelsURL string
}

func (f *FindModel) Request(msg *ma.Message) *ma.Message {
	f.logger.Info("FindModel step is sending request")
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Каталог моделей", WebApp: &tgbotapi.WebAppInfo{
			Url: fmt.Sprintf("%s?city_id=%s&service_id=%s", f.ModelsURL, f.state.GetCityID(), f.state.GetServiceID()),
		}}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: BackToMain}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage("Поиск моделей", msg, keyboard, false)
	}
	return nil
}

func (f *FindModel) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	f.logger.Infof("FindModel step is processing response")
	if Compare(msg.Text, BackToMain) {
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
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: BackToMain}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage("Всем привет! Меня зовут Маша и я алкоголик. Давайте сотрудничать.", msg, keyboard, false)
	}
	return nil
}

func (c *Collaboration) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	c.logger.Infof("Collaboration step is processing response")
	if Compare(msg.Text, BackToMain) {
		c.logger.Infof("Next step is MainMenuStep")
		return nil, PreviousStep
	}
	return nil, EmptyStep
}
