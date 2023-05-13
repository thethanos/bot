package bot

import (
	ma "multimessenger_bot/internal/messenger_adapter"
	"strings"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

type AddService struct {
	StepBase
}

func (a *AddService) Request(msg *ma.Message) *ma.Message {
	a.logger.Info("AddService step is sending request")
	a.inProgress = true
	text := "Введите название услуги"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 1)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Назад"}}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true, OneTimeKeyboard: true}
		return ma.NewMessage(text, ma.REGULAR, msg, keyboard, nil)
	}
	return ma.NewMessage(text, ma.REGULAR, msg, nil, nil)
}

func (a *AddService) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	a.logger.Info("AddService step is processing response")
	a.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		a.logger.Info("Next step is PreviousStep")
		return nil, PreviousStep
	}
	a.DbAdapter.SaveNewService(msg.Text)
	a.logger.Info("Next step is PreviousStep")
	return nil, PreviousStep
}

type AddCity struct {
	StepBase
}

func (a *AddCity) Request(msg *ma.Message) *ma.Message {
	a.logger.Info("AddCity step is sending request")
	a.inProgress = true
	text := "Введите название города"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 1)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Назад"}}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true, OneTimeKeyboard: true}
		return ma.NewMessage(text, ma.REGULAR, msg, keyboard, nil)
	}
	return ma.NewMessage(text, ma.REGULAR, msg, nil, nil)
}

func (a *AddCity) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	a.logger.Info("AddCity step is processing response")
	a.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		a.logger.Info("Next step is PreviousStep")
		return nil, PreviousStep
	}
	a.DbAdapter.SaveNewCity(msg.Text)
	a.logger.Info("Next step is PreviousStep")
	return nil, PreviousStep
}
