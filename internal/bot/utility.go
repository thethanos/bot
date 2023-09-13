package bot

import (
	ma "bot/internal/msgadapter"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

type StepType uint

const (
	MainMenuStep StepType = iota
	MainMenuServiceCategorySelectionStep
	MainMenuServiceSelectionStep
	ServiceCategorySelectionStep
	ServiceSelectionStep
	CitySelectionStep
	MainMenuCitySelectionStep
	MasterSelectionStep
	FindModelStep
	CollaborationStep
	PreviousStep
	EmptyStep
)

func getStepTypeName(step StepType) string {
	switch step {
	case MainMenuStep:
		return "MainMenuStep"
	case MainMenuServiceCategorySelectionStep:
		return "MainMenuServiceCategorySelectionStep"
	case MainMenuServiceSelectionStep:
		return "MainMenuServiceSelectionStep"
	case ServiceCategorySelectionStep:
		return "ServiceCategorySelectionStep"
	case ServiceSelectionStep:
		return "ServiceSelectionStep"
	case CitySelectionStep:
		return "CitySelectionStep"
	case MasterSelectionStep:
		return "MasterSelectionStep"
	case FindModelStep:
		return "FindModelStep"
	case CollaborationStep:
		return "CollaborationStep"
	case EmptyStep:
		return "EmptyStep"
	case PreviousStep:
		return "PreviousStep"
	default:
		return "Unknown type"
	}
}

func makeKeyboard(btnsCreate []string, btnsAppend ...[]tgbotapi.KeyboardButton) *tgbotapi.ReplyKeyboardMarkup {

	rows := make([][]tgbotapi.KeyboardButton, 0)
	for _, text := range btnsCreate {
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: text}})
	}
	if len(rows) > 0 {
		rows = append(rows, btnsAppend...)
	}
	return &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true, OneTimeKeyboard: true}
}

type Question struct {
	Text  string
	Field string
}

type StepStack struct {
	steps []Step
}

func NewStepStack() *StepStack {
	return &StepStack{
		steps: make([]Step, 0),
	}
}

func (s *StepStack) Push(step Step) {
	s.steps = append(s.steps, step)
}

func (s *StepStack) Pop() {
	s.steps = s.steps[:len(s.steps)-1]
}

func (s *StepStack) Top() Step {
	return s.steps[len(s.steps)-1]
}

func (s *StepStack) Empty() bool {
	return len(s.steps) == 0
}

func (s *StepStack) Clear() {
	s.steps = make([]Step, 0)
}

type Step interface {
	ProcessResponse(*ma.Message) (*ma.Message, StepType)
	Request(*ma.Message) *ma.Message
	IsInProgress() bool
	Reset()
	SetInProgress(bool)
}
