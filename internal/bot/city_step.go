package bot

import (
	"bot/internal/dbadapter"
	"bot/internal/entities"
	ma "bot/internal/msgadapter"
	"fmt"
	"strings"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

type CitySelectionStepMode interface {
	MenuItems([]*entities.City) [][]tgbotapi.KeyboardButton
	Buttons() [][]tgbotapi.KeyboardButton
	NextStep() StepType
}

type BaseCitySelectionMode struct {
	dbAdapter *dbadapter.DBAdapter
}

func (b *BaseCitySelectionMode) MenuItems(cities []*entities.City) [][]tgbotapi.KeyboardButton {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	for _, city := range cities {
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: city.Name}})
	}
	return rows
}

func (b *BaseCitySelectionMode) Buttons() [][]tgbotapi.KeyboardButton {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться назад"}})
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться на главную"}})
	return rows
}

func (b *BaseCitySelectionMode) NextStep() StepType {
	return MasterSelectionStep
}

type MainMenuCitySelectionMode struct {
	BaseCitySelectionMode
}

func (m *MainMenuCitySelectionMode) Buttons() [][]tgbotapi.KeyboardButton {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Вернуться на главную"}})
	return rows
}

func (m *MainMenuCitySelectionMode) NextStep() StepType {
	return ServiceCategorySelectionStep
}

type CitySelection struct {
	StepBase
	cities []*entities.City
	mode   CitySelectionStepMode
}

func (c *CitySelection) Request(msg *ma.Message) *ma.Message {
	c.logger.Infof("CitySelection step is sending request")
	c.inProgress = true

	cities, _ := c.DBAdapter.GetCities(c.state.GetServiceID(), 0, -1)

	if msg.Source == ma.TELEGRAM {
		rows := c.mode.MenuItems(cities)
		rows = append(rows, c.mode.Buttons()...)
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(cities) == 0 {
			return ma.NewTextMessage("По вашему запросу ничего не найдено", msg, keyboard, false)
		}

		c.cities = cities
		return ma.NewTextMessage(" Выберите город", msg, keyboard, false)
	}
	return ma.NewTextMessage("this messenger is unsupported yet", msg, nil, true)
}

func (c *CitySelection) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	c.logger.Infof("CitySelection step is processing response")
	c.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "вернуться назад" {
		c.logger.Info("Next step is PreviousStep")
		return nil, PreviousStep
	}
	if userAnswer == "вернуться на главную" {
		c.logger.Info("Next step is MainMenuStep")
		return nil, MainMenuStep
	}

	for idx, city := range c.cities {
		if userAnswer == strings.ToLower(city.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			c.state.City = city
			c.logger.Infof("Next step is %s", getStepTypeName(c.mode.NextStep()))
			return nil, c.mode.NextStep()
		}
	}

	c.inProgress = true
	c.logger.Info("Next step is EmptyStep")
	return ma.NewTextMessage("Пожалуйста выберите ответ из списка.", msg, nil, false), EmptyStep
}

func (c *CitySelection) Reset() {
	c.state.City = nil
}
