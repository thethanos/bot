package bot

import (
	"fmt"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/entities"
	ma "multimessenger_bot/internal/messenger_adapter"
	"strings"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
	"go.uber.org/zap"
)

type StepType uint

const (
	MainMenuStep StepType = iota
	MainMenuRequestStep
	MainMenuServiceCategorySelectionStep
	MainMenuServiceSelectionStep
	ServiceCategorySelectionStep
	ServiceSelectionStep
	CitySelectionStep
	MainMenuCitySelectionStep
	MasterSelectionStep
	AboutStep
	MasterStep
	MasterRegistrationStep
	MasterCityPromptStep
	MasterServiceCategorySecletionStep
	MasterServiceSelectionStep
	MasterRegistrationFinalStep
	PreviousStep
	AdminStep
	AdminServiceCategorySelectionStep
	AddServiceCategoryStep
	AddServiceStep
	AddCityStep
	AddMasterStep
	AddMasterFinalStep
	ImageUploadStep
	EmptyStep
)

func getStepTypeName(step StepType) string {
	switch step {
	case MainMenuStep:
		return "MainMenuStep"
	case MainMenuRequestStep:
		return "MainMenuRequestStep"
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
	case AboutStep:
		return "AboutStep"
	case MasterStep:
		return "MasterStep"
	case EmptyStep:
		return "EmptyStep"
	case MasterRegistrationStep:
		return "MasterRegistrationStep"
	case MasterServiceCategorySecletionStep:
		return "MasterServiceCategorySecletionStep"
	case MasterServiceSelectionStep:
		return "MasterServiceSelectionStep"
	case MasterRegistrationFinalStep:
		return "RegistrationFinalStep"
	case PreviousStep:
		return "PreviousStep"
	case AdminStep:
		return "AdminStep"
	case AddServiceCategoryStep:
		return "AddServiceCategoryStep"
	case AddServiceStep:
		return "AddServiceStep"
	case AddCityStep:
		return "AddCityStep"
	case AddMasterStep:
		return "AddMasterStep"
	case AddMasterFinalStep:
		return "AddMasterFinalStep"
	case ImageUploadStep:
		return "ImageUploadStep"
	default:
		return "Unknown type"
	}
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

type StepBase struct {
	logger     *zap.SugaredLogger
	inProgress bool
	state      *entities.UserState
	dbAdapter  *db_adapter.DbAdapter
}

func (s *StepBase) IsInProgress() bool {
	return s.inProgress
}

func (s *StepBase) Reset() {
}

func (s *StepBase) SetInProgress(flag bool) {
	s.inProgress = flag
}

type YesNo struct {
	StepBase
	question Question
	yesStep  StepType
	noStep   StepType
}

func (y *YesNo) Request(msg *ma.Message) *ma.Message {
	y.logger.Infof("YesNo step is sending request")
	y.inProgress = true
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 2)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Да"}}
		rows[1] = []tgbotapi.KeyboardButton{{Text: "Нет"}}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true, OneTimeKeyboard: true}
		return ma.NewTextMessage(y.question.Text, msg, keyboard, false)
	}
	return ma.NewTextMessage(fmt.Sprintf("%s\n1. Да\n2. Нет", y.question.Text), msg, nil, true)
}

func (y *YesNo) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	y.logger.Infof("YesNo step is processing response")
	y.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "да" || userAnswer == "1" {
		y.logger.Infof("Next step is %s", getStepTypeName(y.yesStep))
		return nil, y.yesStep
	}
	y.logger.Infof("Next step is %s", getStepTypeName(y.yesStep))
	return nil, y.noStep
}

type Prompt struct {
	StepBase
	question Question
	nextStep StepType
	errStep  StepType
}

func (p *Prompt) Request(msg *ma.Message) *ma.Message {
	p.logger.Infof("Prompt step is sending request")
	p.inProgress = true
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 1)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Назад"}}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true, OneTimeKeyboard: true}
		return ma.NewTextMessage(p.question.Text, msg, keyboard, false)
	}

	return ma.NewTextMessage(p.question.Text, msg, nil, true)
}

func (p *Prompt) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	p.logger.Infof("Prompt step is processing response")
	p.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		p.logger.Info("Next step is PreviousStep")
		return nil, PreviousStep
	}
	p.state.RawInput[p.question.Field] = msg.Text
	p.logger.Infof("Next step is %s", getStepTypeName(p.nextStep))
	return nil, p.nextStep
}
