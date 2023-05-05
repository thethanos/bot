package bot

import (
	"fmt"
	ci "multimessenger_bot/internal/client_interface"
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
	ServicesStep
	CitiesStep
	QuestionsStep
	AboutStep
	MasterStep
	FinalStep
	EmptyStep
)

type Step interface {
	ProcessResponse(ci.Message) (ci.Message, int)
	Request(ci.Message) ci.Message
	DefaultRequest(ci.Message) ci.Message
	IsInProgress() bool
}

type MainMenu struct {
	inProgress bool
	State      *UserState
}

func (m *MainMenu) Request(msg ci.Message) ci.Message {
	text := "1) услуги\n2) город\n3) вопросы\n4) о нас\n5)мастер"
	m.inProgress = true
	return ci.Message{Text: text, WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}
}

func (m *MainMenu) DefaultRequest(msg ci.Message) ci.Message {
	return ci.Message{Text: "Хотите вернуться в главное меню?\nДа", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}
}

func (m *MainMenu) ProcessResponse(msg ci.Message) (ci.Message, int) {
	m.inProgress = false

	switch msg.Text {
	case "услуги":
		return ci.Message{WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, ServicesStep
	case "город":
		return ci.Message{WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, CitiesStep
	case "вопросы":
		return ci.Message{WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, QuestionsStep
	case "о нас":
		return ci.Message{WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, AboutStep
	case "мастер":
		return ci.Message{WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, MasterStep
	case "Да":
		return ci.Message{WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, MainMenuStep
	}

	return ci.Message{Text: "Пожалуйста выберите ответ из списка.", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, EmptyStep
}

func (m *MainMenu) IsInProgress() bool {
	return m.inProgress
}

type Cities struct {
	inProgress bool
	State      *UserState
}

func (c *Cities) Request(msg ci.Message) ci.Message {
	text := "1) Тель-Авив\n 2) Нетания\n"
	c.inProgress = true
	return ci.Message{Text: text, WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}
}

func (c *Cities) DefaultRequest(msg ci.Message) ci.Message {
	c.inProgress = false
	return ci.Message{Text: "something went wrong", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}
}

func (c *Cities) ProcessResponse(msg ci.Message) (ci.Message, int) {
	c.inProgress = false

	switch msg.Text {
	case "Тель-Авив", "Нетания":
		c.State.city = msg.Text
	default:
		c.inProgress = true
		return ci.Message{Text: "Пожалуйста выберите ответ из списка.", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, EmptyStep
	}

	if len(c.State.service) == 0 {
		return ci.Message{Text: "", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, ServicesStep
	} else {
		return ci.Message{Text: "", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, MasterStep
	}
}

func (c *Cities) IsInProgress() bool {
	return c.inProgress
}

type Services struct {
	inProgress bool
	State      *UserState
}

func (c *Services) Request(msg ci.Message) ci.Message {
	text := "1) услуга1\n 2) услуга2\n"
	c.inProgress = true
	return ci.Message{Text: text, WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}
}

func (c *Services) DefaultRequest(msg ci.Message) ci.Message {
	c.inProgress = false
	return ci.Message{Text: "something went wrong", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}
}

func (c *Services) ProcessResponse(msg ci.Message) (ci.Message, int) {
	c.inProgress = false

	switch msg.Text {
	case "услуга1", "услуга2":
		c.State.service = msg.Text
	default:
		c.inProgress = true
		return ci.Message{Text: "Пожалуйста выберите ответ из списка.", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, EmptyStep
	}

	if len(c.State.city) == 0 {
		return ci.Message{Text: "", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, CitiesStep
	} else {
		return ci.Message{Text: "", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, MasterStep
	}
}

func (c *Services) IsInProgress() bool {
	return c.inProgress
}

type Master struct {
	inProgress bool
	State      *UserState
}

func (m *Master) Request(msg ci.Message) ci.Message {
	text := "1) мастер1\n 2) мастер2\n"
	m.inProgress = true
	return ci.Message{Text: text, WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}
}

func (m *Master) DefaultRequest(msg ci.Message) ci.Message {
	m.inProgress = false
	return ci.Message{Text: "something went wrong", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}
}

func (m *Master) ProcessResponse(msg ci.Message) (ci.Message, int) {
	m.inProgress = false

	switch msg.Text {
	case "мастер1", "мастер2":
		m.State.master = msg.Text
		return ci.Message{WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, FinalStep
	default:
		m.inProgress = true
		return ci.Message{Text: "Пожалуйста выберите ответ из списка.", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, EmptyStep
	}
}

func (m *Master) IsInProgress() bool {
	return m.inProgress
}

type Final struct {
	inProgress bool
	State      *UserState
}

func (f *Final) Request(msg ci.Message) ci.Message {
	text := fmt.Sprintf("Ваша запись\nУслуга: %s\nГород: %s\nМастер: %s\nПодтвердить?\nДа\nНет", f.State.service, f.State.city, f.State.master)
	f.inProgress = true
	return ci.Message{Text: text, WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}
}

func (f *Final) DefaultRequest(msg ci.Message) ci.Message {
	f.inProgress = false
	return ci.Message{Text: "something went wrong", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}
}

func (f *Final) ProcessResponse(msg ci.Message) (ci.Message, int) {
	f.inProgress = false

	switch msg.Text {
	case "Да":
		return ci.Message{Text: "Запись завершена", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, MainMenuStep
	case "Нет":
		return ci.Message{Text: "Запись отменена", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, MainMenuStep
	default:
		f.inProgress = true
		return ci.Message{Text: "Пожалуйста выберите ответ из списка.", WaData: msg.WaData, TgData: msg.TgData, Type: msg.Type}, EmptyStep
	}
}

func (f *Final) IsInProgress() bool {
	return f.inProgress
}

type Empty struct {
}

func (e *Empty) Request(msg ci.Message) ci.Message {
	return msg
}

func (e *Empty) DefaultRequest(msg ci.Message) ci.Message {
	return msg
}

func (e *Empty) ProcessResponse(msg ci.Message) (ci.Message, int) {
	return msg, EmptyStep
}

func (e *Empty) IsInProgress() bool {
	return false
}
