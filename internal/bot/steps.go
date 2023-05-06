package bot

import (
	"fmt"
	ma "multimessenger_bot/internal/messenger_adapter"
)

type UserSession struct {
	CurrentStep Step
	State       UserState
}

type UserState struct {
	city    string
	service string
	master  string
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
}

func (c *CitySelection) Request(msg *ma.Message) *ma.Message {
	text := "1) Тель-Авив\n2) Нетания\n"
	c.inProgress = true
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (c *CitySelection) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	c.inProgress = false

	switch msg.Text {
	case "Тель-Авив", "Нетания":
		c.State.city = msg.Text
	default:
		c.inProgress = true
		return &ma.Message{Text: "Пожалуйста выберите ответ из списка.", UserData: msg.UserData, Type: msg.Type}, EmptyStep
	}

	if len(c.State.service) == 0 {
		return nil, ServiceSelectionStep
	} else {
		return nil, MasterSelectionStep
	}
}

func (c *CitySelection) IsInProgress() bool {
	return c.inProgress
}

type ServiceSelection struct {
	StepBase
}

func (c *ServiceSelection) Request(msg *ma.Message) *ma.Message {
	text := "1) услуга1\n2) услуга2\n"
	c.inProgress = true
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (c *ServiceSelection) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	c.inProgress = false

	switch msg.Text {
	case "услуга1", "услуга2":
		c.State.service = msg.Text
	default:
		c.inProgress = true
		return &ma.Message{Text: "Пожалуйста выберите ответ из списка.", UserData: msg.UserData, Type: msg.Type}, EmptyStep
	}

	if len(c.State.city) == 0 {
		return nil, CitySelectionStep
	} else {
		return nil, MasterSelectionStep
	}
}

func (c *ServiceSelection) IsInProgress() bool {
	return c.inProgress
}

type MasterSelection struct {
	StepBase
}

func (m *MasterSelection) Request(msg *ma.Message) *ma.Message {
	text := "1) мастер1\n2) мастер2\n"
	m.inProgress = true
	return &ma.Message{Text: text, UserData: msg.UserData, Type: msg.Type}
}

func (m *MasterSelection) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	m.inProgress = false

	switch msg.Text {
	case "мастер1", "мастер2":
		m.State.master = msg.Text
		return nil, FinalStep
	default:
		m.inProgress = true
		return &ma.Message{Text: "Пожалуйста выберите ответ из списка.", UserData: msg.UserData, Type: msg.Type}, EmptyStep
	}
}

func (m *MasterSelection) IsInProgress() bool {
	return m.inProgress
}

type Final struct {
	StepBase
}

func (f *Final) Request(msg *ma.Message) *ma.Message {
	text := fmt.Sprintf("Ваша запись\nУслуга: %s\nГород: %s\nМастер: %s\nПодтвердить?\nДа\nНет", f.State.service, f.State.city, f.State.master)
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
