package bot

import (
	"fmt"
	"multimessenger_bot/internal/entities"
	ma "multimessenger_bot/internal/messenger_adapter"
	"strings"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

type CityPromptStepMode interface {
	Text() string
	Buttons() *tgbotapi.ReplyKeyboardMarkup
	NextStep() StepType
}

type BaseCityPromptMode struct {
}

func (b *BaseCityPromptMode) Text() string {
	return "Введите город"
}

func (b *BaseCityPromptMode) Buttons() *tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
	return &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
}

func (b *BaseCityPromptMode) NextStep() StepType {
	return EmptyStep
}

type MainMenuCityPromptMode struct {
	BaseCityPromptMode
}

func (m *MainMenuCityPromptMode) Buttons() *tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
	return &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
}

func (m *MainMenuCityPromptMode) NextStep() StepType {
	return ServiceCategorySelectionStep
}

type RegistrationCityPromptMode struct {
	BaseCityPromptMode
}

func (r *RegistrationCityPromptMode) NextStep() StepType {
	return MasterServiceCategorySecletionStep
}

type CityPrompt struct {
	StepBase
	mode CityPromptStepMode
}

func (c *CityPrompt) Request(msg *ma.Message) *ma.Message {
	c.logger.Infof("CityPrompt step is sending request")
	c.inProgress = true

	if msg.Source == ma.TELEGRAM {
		return ma.NewMessage(c.mode.Text(), ma.REGULAR, msg, c.mode.Buttons(), nil)
	}

	return ma.NewMessage(fmt.Sprintf("%s\n1. Назад\n2. Главное меню", c.mode.Text()), ma.REGULAR, msg, nil, nil)
}

func (c *CityPrompt) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	c.logger.Infof("CityPrompt step is processing response")
	c.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		return nil, PreviousStep
	}
	if userAnswer == "главное меню" {
		return nil, MainMenuStep
	}

	city, err := c.dbAdapter.GetCity(msg.Text)
	if err != nil {
		c.inProgress = true
		c.logger.Infof("Next step is CityPromptStep")
		return ma.NewMessage(fmt.Sprintf("По запросу %s ничего не найдено", msg.Text), ma.REGULAR, msg, nil, nil), CityPromptStep
	}
	c.state.City = city
	c.logger.Infof("Next step is %s", getStepTypeName(c.mode.NextStep()))
	return nil, c.mode.NextStep()
}

func (c *CityPrompt) Reset() {
	c.state.City = nil
}

type CitySelection struct {
	StepBase
	cities       []*entities.City
	filter       bool
	checkService bool
	nextStep     StepType
	errStep      StepType
}

func (c *CitySelection) Request(msg *ma.Message) *ma.Message {
	c.logger.Infof("CitySelection step is sending request")
	c.inProgress = true

	cities, _ := c.dbAdapter.GetCities(c.state.Service.ID)

	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		for _, city := range cities {
			rows = append(rows, []tgbotapi.KeyboardButton{{Text: city.Name}})
		}
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(cities) == 0 {
			return ma.NewMessage("По вашему запросу ничего не найдено", ma.REGULAR, msg, keyboard, nil)
		}

		c.cities = cities
		return ma.NewMessage(" Выберите город", ma.REGULAR, msg, keyboard, nil)
	}

	text := ""
	for idx, city := range cities {
		text += fmt.Sprintf("%d. %s\n", idx+1, city.Name)
	}
	text += fmt.Sprintf("%d. Назад\n", len(cities)+1)

	c.cities = cities
	return ma.NewMessage(text, ma.REGULAR, msg, nil, nil)
}

func (c *CitySelection) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	c.logger.Infof("CitySelection step is processing response")
	if msg.Type == ma.CALLBACK {
		return nil, EmptyStep
	}
	c.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		return nil, PreviousStep
	}
	if userAnswer == "главное меню" {
		return nil, MainMenuStep
	}

	for idx, city := range c.cities {
		if userAnswer == strings.ToLower(city.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			c.state.City = city
			return nil, EmptyStep
		}
	}

	c.inProgress = true
	return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), c.errStep
}

func (c *CitySelection) Reset() {
	c.state.City = nil
}
