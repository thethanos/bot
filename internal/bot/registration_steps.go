package bot

import (
	"fmt"
	ma "multimessenger_bot/internal/messenger_adapter"
	"strings"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

type Question struct {
	Text  string
	Field string
}

type RegistrationFinal struct {
	StepBase
}

func (r *RegistrationFinal) Request(msg *ma.Message) *ma.Message {
	r.logger.Infof("RegistrationFinal step is sending request")
	r.inProgress = true
	data := FormatMapToString(r.state.RawInput)
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Да"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Нет"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage(fmt.Sprintf("%s\nПодтвердить регистрацию?", data), msg, keyboard, false)
	}
	return ma.NewTextMessage(fmt.Sprintf("%s\nПодтвердить регистрацию?\n1. Да\n2. Нет", data), msg, nil, true)
}

func (r *RegistrationFinal) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	r.logger.Infof("RegistrationFinal step is processing response")
	r.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "да" || userAnswer == "1" {
		r.dbAdapter.SaveMaster(r.state)
		r.state.Reset()
		return ma.NewTextMessage("Регистрация прошла успешно!", msg, nil, true), MainMenuRequestStep
	}
	r.state.Reset()
	return nil, MainMenuRequestStep
}

func FormatMapToString(data map[string]string) string {
	res := ""
	for key, val := range data {
		res += fmt.Sprintf("%s: %s\n", key, val)
	}
	return res
}
