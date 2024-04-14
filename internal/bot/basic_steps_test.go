package bot

import (
	"bot/internal/config"
	"bot/internal/entities"
	"bot/internal/logger"
	ma "bot/internal/msgadapter"
	"reflect"
	"testing"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

const (
	unsupported = "this messenger is unsupported yet"
)

func TestYesNoStep(t *testing.T) {

	text := "test text"
	step := &YesNo{
		StepBase: StepBase{
			logger: logger.NewLogger(config.RELEASE),
		},
		question: Question{
			Text: text,
		},
		yesStep: MainMenuCitySelectionStep,
		noStep:  MainMenuServiceSelectionStep,
	}

	msg := &ma.Message{
		Text:   text,
		Source: ma.TELEGRAM,
		Data: &ma.MessageData{
			TgMarkup:     makeKeyboard([]string{"Да", "Нет"}),
			RemoveMarkup: false,
		},
	}

	if res := step.Request(msg); !reflect.DeepEqual(res, msg) {
		t.Error("YesNo step returned wrong message")
	}

	msg.Source = ma.WHATSAPP
	if res := step.Request(msg); res.Text != unsupported {
		t.Error("YesNo step returned wrong message")
	}

	if res, _ := step.ProcessResponse(msg); res != nil {
		t.Error("YesNo step ProcessResponse returned not nil message")
	}

	resp := ma.NewTextMessage("Да", msg, nil, true)
	if _, nextStep := step.ProcessResponse(resp); nextStep != MainMenuCitySelectionStep {
		t.Error("YesNo step returned wrong next step")
	}

	resp = ma.NewTextMessage("Нет", msg, nil, true)
	if _, nextStep := step.ProcessResponse(resp); nextStep != MainMenuServiceSelectionStep {
		t.Error("YesNo step returned wrong next step")
	}
}

func TestPromptStep(t *testing.T) {

	text := "test text"
	step := &Prompt{
		StepBase: StepBase{
			logger: logger.NewLogger(config.RELEASE),
			state: &entities.UserState{
				RawInput: make(map[string]string),
			},
		},
		question: Question{
			Text: text,
		},
		nextStep: MainMenuCitySelectionStep,
	}

	msg := &ma.Message{
		Text:   text,
		Source: ma.TELEGRAM,
		Data: &ma.MessageData{
			TgMarkup:     makeKeyboard([]string{Back}),
			RemoveMarkup: false,
		},
	}

	if res := step.Request(msg); !reflect.DeepEqual(res, msg) {
		t.Error("Prompt step returned wrong message")
	}

	msg.Source = ma.WHATSAPP
	if res := step.Request(msg); res.Text != unsupported {
		t.Error("Prompt step returned wrong message")
	}

	if res, _ := step.ProcessResponse(msg); res != nil {
		t.Error("Prompt step ProcessResponse returned not nil message")
	}

	resp := ma.NewTextMessage(Back, msg, nil, true)
	if _, nextStep := step.ProcessResponse(resp); nextStep != PreviousStep {
		t.Error("Prompt step returned wrong next step")
	}

	resp = ma.NewTextMessage(text, msg, nil, true)
	if _, nextStep := step.ProcessResponse(resp); nextStep != MainMenuCitySelectionStep {
		t.Error("Prompt step returned wrong next step")
	}
}

func TestMainMenuStep(t *testing.T) {

	step := &MainMenu{
		StepBase: StepBase{
			logger: logger.NewLogger(config.RELEASE),
			state:  &entities.UserState{},
		},
	}

	text := "Главное меню"
	msg := &ma.Message{
		Text:   text,
		Source: ma.TELEGRAM,
		Data: &ma.MessageData{
			TgMarkup:     makeKeyboard([]string{"Город", "Услуги", "Поиск моделей", "По вопросам сотрудничества"}),
			RemoveMarkup: false,
		},
	}

	msg.Source = ma.WHATSAPP
	if res := step.Request(msg); res.Text != unsupported {
		t.Error("MainMenu step returned wrong message")
	}

	msg = ma.NewTextMessage("Пожалуйста выберите ответ из списка.", msg, nil, false)
	if res, nextStep := step.ProcessResponse(msg); !reflect.DeepEqual(res, msg) || nextStep != EmptyStep {
		t.Error("MainMenu step ProcessResponse returned wrong message")
	}

	resp := ma.NewTextMessage("Город", msg, nil, true)
	if _, nextStep := step.ProcessResponse(resp); nextStep != MainMenuCitySelectionStep {
		t.Error("MainMenu step returned wrong next step")
	}

	resp = ma.NewTextMessage("Услуги", msg, nil, true)
	if _, nextStep := step.ProcessResponse(resp); nextStep != MainMenuServiceCategorySelectionStep {
		t.Error("MainMenu step returned wrong next step")
	}

	resp = ma.NewTextMessage("Поиск моделей", msg, nil, true)
	if _, nextStep := step.ProcessResponse(resp); nextStep != FindModelStep {
		t.Error("MainMenu step returned wrong next step")
	}

	resp = ma.NewTextMessage("По вопросам сотрудничества", msg, nil, true)
	if _, nextStep := step.ProcessResponse(resp); nextStep != CollaborationStep {
		t.Error("MainMenu step returned wrong next step")
	}
}

func TestMasterSelectionStep(t *testing.T) {

	step := MasterSelection{
		StepBase: StepBase{
			logger: logger.NewLogger(config.RELEASE),
			state: &entities.UserState{
				City: &entities.City{
					ID: 0,
				},
				Service: &entities.Service{
					ID: 0,
				},
			},
		},
	}

	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Каталог мастеров", WebApp: &tgbotapi.WebAppInfo{
		Url: "https://bot-dev-domain.com:1445/bot-webapp/gallery?city_id=123&service_id=123",
	}}})
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: Back}})
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: BackToMain}})
	keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

	text := "Выбор мастера"
	msg := &ma.Message{
		Text:   text,
		Source: ma.TELEGRAM,
		UserID: "test",
		Data: &ma.MessageData{
			TgMarkup:     keyboard,
			RemoveMarkup: false,
		},
	}

	msg.Source = ma.WHATSAPP
	if res := step.Request(msg); res.Text != unsupported {
		t.Error("MasterSelection step returned wrong message")
	}

	resp := ma.NewTextMessage(Back, msg, nil, true)
	if _, nextStep := step.ProcessResponse(resp); nextStep != PreviousStep {
		t.Error("MasterSelection step returned wrong next step")
	}

	resp = ma.NewTextMessage(BackToMain, msg, nil, true)
	if _, nextStep := step.ProcessResponse(resp); nextStep != MainMenuStep {
		t.Error("MasterSelection step returned wrong next step")
	}
}
