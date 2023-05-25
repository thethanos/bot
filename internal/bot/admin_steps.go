package bot

import (
	ma "multimessenger_bot/internal/messenger_adapter"
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
	a.dbAdapter.SaveNewServiceCategory(msg.Text)
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
		rows := make([][]tgbotapi.KeyboardButton, 1)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Назад"}}
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
	a.dbAdapter.SaveNewService(msg.Text, a.state.ServiceCategory.ID)
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
	a.dbAdapter.SaveNewCity(msg.Text)
	a.logger.Info("Next step is PreviousStep")
	return nil, PreviousStep
}

type AddMaster struct {
	StepBase
}

func (a *AddMaster) Request(msg *ma.Message) *ma.Message {
	a.logger.Info("AddMaster step is sending request")
	a.inProgress = true
	return nil
}

func (a *AddMaster) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	a.logger.Info("AddMaster step is processing response")
	a.inProgress = false
	return nil, ImageUploadStep
}

type Downloader interface {
	DownloadFile(msg *ma.Message)
}

type ImageUpload struct {
	StepBase
	downloader Downloader
}

func (i *ImageUpload) Request(msg *ma.Message) *ma.Message {
	i.logger.Info("ImageUpload step is sending request")
	i.inProgress = true
	return nil
}

func (i *ImageUpload) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	i.logger.Info("ImageUpload step is processing response")
	i.downloader.DownloadFile(msg)
	return nil, EmptyStep
}
