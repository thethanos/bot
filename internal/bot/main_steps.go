package bot

import (
	ma "multimessenger_bot/internal/messenger_adapter"
	"strings"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

type Greetings struct {
	StepBase
}

func (g *Greetings) Request(msg *ma.Message) *ma.Message {
	return nil
}

func (g *Greetings) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	return nil, EmptyStep
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
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "О нас"}})
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
		return nil, MainMenuCitySelectionStep
	case "услуги":
		return nil, MainMenuServiceCategorySelectionStep
	case "поиск моделей":
		return nil, EmptyStep
	case "по вопросам сотрудничества":
		return nil, MasterStep
	case "о нас":
		return nil, AboutStep
	case "админ":
		return nil, AdminStep
	}

	return ma.NewTextMessage("Пожалуйста выберите ответ из списка.", msg, nil, false), EmptyStep
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
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Добавить мастера"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться на главную"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage(text, msg, keyboard, false)
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
	case "добавить мастера":
		a.logger.Info("Next step is AddMaster")
		return nil, AddMasterStep
	default:
		a.inProgress = true
		a.logger.Info("Next step is EmptyStep")
		return ma.NewTextMessage("Пожалуйста выберите ответ из списка.", msg, nil, false), EmptyStep
	}
}
