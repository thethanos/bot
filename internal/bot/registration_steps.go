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

var MasterQuestions = []*Question{
	{Text: "Как вас называть?", Field: "name"},
	{Text: "В каком городе вы работаете?", Field: "city"},
	{Text: "Какую услугу предоставляете?", Field: "service"},
}

type RegistrationFinal struct {
	StepBase
}

func (r *RegistrationFinal) Request(msg *ma.Message) *ma.Message {
	r.inProgress = true
	data := FormatMapToString(r.State.RawInput)
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 2)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Да"}}
		rows[1] = []tgbotapi.KeyboardButton{{Text: "Нет"}}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewMessage(fmt.Sprintf("%s\nПодтвердить регистрацию?", data), ma.REGULAR, msg, keyboard, nil)
	}
	return ma.NewMessage(fmt.Sprintf("%s\nПодтвердить регистрацию?\n1. Да\n2. Нет", data), ma.REGULAR, msg, nil, nil)
}

func (r *RegistrationFinal) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {

	if msg.Type == ma.CALLBACK {
		return nil, EmptyStep
	}
	r.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "да" || userAnswer == "1" {
		r.DbAdapter.SaveNewMaster(r.State)
		r.State.Reset()
		return ma.NewMessage("Регистрация прошла успешно!", ma.REGULAR, msg, nil, nil), MainMenuRequestStep
	}
	r.State.Reset()
	return nil, MainMenuRequestStep
}

func FormatMapToString(data map[string]string) string {
	res := ""
	for key, val := range data {
		res += fmt.Sprintf("%s: %s\n", key, val)
	}
	return res
}
