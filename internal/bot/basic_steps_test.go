package bot

import (
	"multimessenger_bot/internal/config"
	"multimessenger_bot/internal/entities"
	"multimessenger_bot/internal/logger"
	ma "multimessenger_bot/internal/messenger_adapter"
	"reflect"
	"testing"
)

const unsupported = "this messenger is unsupported yet"

func TestStepBase(t *testing.T) {

	base := StepBase{}
	base.SetInProgress(true)
	if base.IsInProgress() != true {
		t.Error("Step is not in progress")
	}

	base.SetInProgress(false)
	if base.IsInProgress() != false {
		t.Error("Step is in progress")
	}
}

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

	if step.IsInProgress() != true {
		t.Error("YesNo step is not in progress after sending request")
	}

	if res, _ := step.ProcessResponse(msg); res != nil {
		t.Error("YesNo step ProcessResponse returned not nil message")
	}

	if step.IsInProgress() != false {
		t.Error("YesNo step is in progress after processing response")
	}

	msg.Source = ma.TELEGRAM
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
			TgMarkup:     makeKeyboard([]string{"Назад"}),
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

	if step.IsInProgress() != true {
		t.Error("Prompt step is not in progress after sending request")
	}

	if res, _ := step.ProcessResponse(msg); res != nil {
		t.Error("Prompt step ProcessResponse returned not nil message")
	}

	if step.IsInProgress() != false {
		t.Error("Prompt step is in progress after processing response")
	}

	msg.Source = ma.TELEGRAM
	resp := ma.NewTextMessage("Назад", msg, nil, true)
	if _, nextStep := step.ProcessResponse(resp); nextStep != PreviousStep {
		t.Error("Prompt step returned wrong next step")
	}

	resp = ma.NewTextMessage(text, msg, nil, true)
	if _, nextStep := step.ProcessResponse(resp); nextStep != MainMenuCitySelectionStep {
		t.Error("Prompt step returned wrong next step")
	}
}
