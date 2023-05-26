package bot

import (
	"fmt"
	ma "multimessenger_bot/internal/messenger_adapter"
	"multimessenger_bot/internal/parsers"
	"strings"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

type AddServiceCategory struct {
	StepBase
}

func (a *AddServiceCategory) Request(msg *ma.Message) *ma.Message {
	a.logger.Info("AddServiceCategory step is sending request")
	a.inProgress = true
	text := "Введите название категории услуги"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 1)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Назад"}}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage(text, msg, keyboard, false)
	}
	return ma.NewTextMessage(text, msg, nil, true)
}

func (a *AddServiceCategory) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	a.logger.Info("AddServiceCategory step is processing response")
	a.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		a.logger.Info("Next step is PreviousStep")
		return nil, PreviousStep
	}
	a.dbAdapter.SaveServiceCategory(msg.Text)
	a.logger.Info("Next step is PreviousStep")
	return nil, PreviousStep
}

type AddService struct {
	StepBase
}

func (a *AddService) Request(msg *ma.Message) *ma.Message {
	a.logger.Info("AddService step is sending request")
	a.inProgress = true
	text := "Введите название услуги"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage(text, msg, keyboard, false)
	}
	return ma.NewTextMessage(text, msg, nil, true)
}

func (a *AddService) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	a.logger.Info("AddService step is processing response")
	a.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		a.logger.Info("Next step is PreviousStep")
		return nil, PreviousStep
	}
	a.dbAdapter.SaveService(msg.Text, a.state.ServiceCategory.ID)
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
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage(text, msg, keyboard, false)
	}
	return ma.NewTextMessage(text, msg, nil, true)
}

func (a *AddCity) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	a.logger.Info("AddCity step is processing response")
	a.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		a.logger.Info("Next step is PreviousStep")
		return nil, PreviousStep
	}
	a.dbAdapter.SaveCity(msg.Text)
	a.logger.Info("Next step is PreviousStep")
	return nil, PreviousStep
}

type AddMaster struct {
	StepBase
}

func (a *AddMaster) Request(msg *ma.Message) *ma.Message {
	a.logger.Info("AddMaster step is sending request")
	a.inProgress = true
	text := "Введите данные мастера"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage(text, msg, keyboard, false)
	}
	return ma.NewTextMessage(text, msg, nil, true)
}

func (a *AddMaster) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	a.logger.Info("AddMaster step is processing response")
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		return nil, PreviousStep
	}
	if userAnswer == "главное меню" {
		return nil, MainMenuStep
	}
	master, err := parsers.ParseMasterData(msg.Text)
	if err != nil {
		return ma.NewTextMessage("Не удалось распарсить данные мастера, проверьте правильность ввода и попробуйте еще раз", msg, nil, false), EmptyStep
	}
	a.state.Master = master
	a.inProgress = false
	return nil, ImageUploadStep
}

type Downloader interface {
	DownloadFile(id string, msg *ma.Message) string
}

type ImageUpload struct {
	StepBase
	downloader Downloader
}

func (i *ImageUpload) Request(msg *ma.Message) *ma.Message {
	i.logger.Info("ImageUpload step is sending request")
	i.inProgress = true
	text := "Добавить фото"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Далее"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage(text, msg, keyboard, false)
	}
	return ma.NewTextMessage(text, msg, nil, false)
}

func (i *ImageUpload) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	i.logger.Info("ImageUpload step is processing response")
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "далее" {
		if err := i.dbAdapter.SaveMasterPreview(i.state.Master); err != nil {
			return ma.NewTextMessage("Не удалось сохранить данные для предпросмотра", msg, nil, false), EmptyStep
		}
		return nil, AddMasterFinalStep
	}
	if userAnswer == "назад" {
		return nil, PreviousStep
	}
	if userAnswer == "главное меню" {
		return nil, MainMenuStep
	}
	i.state.Master.Images = append(i.state.Master.Images, i.downloader.DownloadFile(i.state.Master.ID, msg))
	return nil, EmptyStep
}

type AddMasterFinal struct {
	StepBase
}

func (a *AddMasterFinal) Request(msg *ma.Message) *ma.Message {
	a.logger.Info("AddMasterFinal step is sending request")
	a.inProgress = true
	text := "Завершающий этап"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Предпросмотр", WebApp: &tgbotapi.WebAppInfo{Url: fmt.Sprintf("https://bot-dev-domain.com/master/preview?master=%s", a.state.Master.ID)}}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Сохранить анкету"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage(text, msg, keyboard, false)
	}
	return ma.NewTextMessage(text, msg, nil, false)
}

func (a *AddMasterFinal) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	a.logger.Info("AddMasterFinal step is processing response")
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		return nil, PreviousStep
	}
	if userAnswer == "главное меню" {
		return nil, MainMenuStep
	}
	return ma.NewTextMessage("Анкета сохранена", msg, nil, false), AdminStep
}
