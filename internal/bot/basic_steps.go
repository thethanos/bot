package bot

import (
	"fmt"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/entities"
	ma "multimessenger_bot/internal/messenger_adapter"
	"strings"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
)

type StepType uint

const (
	MainMenuStep StepType = iota
	MainMenuRequestStep
	ServiceSelectionStep
	CitySelectionStep
	QuestionsStep
	AboutStep
	MasterSelectionStep
	MasterStep
	FinalStep
	EmptyStep
	RegistrationStep
	RegistrationStepService
	RegistrationStepCity
	RegistrationFinalStep
	PreviousStep
	TestStep
)

type StepStack struct {
	steps []Step
}

func (s *StepStack) Push(step Step) {
	if s.steps == nil {
		s.steps = make([]Step, 0)
	}
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

type Step interface {
	ProcessResponse(*ma.Message) (*ma.Message, StepType)
	Request(*ma.Message) *ma.Message
	IsInProgress() bool
	IsCallBackStep() bool
	Reset()
	SetInProgress(bool)
}

type StepBase struct {
	inProgress bool
	State      *entities.UserState
	DbAdapter  *db_adapter.DbAdapter
}

func (s *StepBase) IsInProgress() bool {
	return s.inProgress
}

func (s *StepBase) IsCallBackStep() bool {
	return false
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
	y.inProgress = true
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 2)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Да"}}
		rows[1] = []tgbotapi.KeyboardButton{{Text: "Нет"}}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true, OneTimeKeyboard: true}
		return ma.NewMessage(y.question.Text, ma.REGULAR, msg, keyboard, nil)
	}
	return ma.NewMessage(fmt.Sprintf("%s\n1. Да\n2. Нет", y.question.Text), ma.REGULAR, msg, nil, nil)
}

func (y *YesNo) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {

	if msg.Type == ma.CALLBACK {
		return nil, EmptyStep
	}

	y.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "да" || userAnswer == "1" {
		return nil, y.yesStep
	}
	return nil, y.noStep
}

type Prompt struct {
	StepBase
	question Question
	nextStep StepType
	errStep  StepType
}

func (p *Prompt) Request(msg *ma.Message) *ma.Message {
	p.inProgress = true
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 1)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Назад"}}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true, OneTimeKeyboard: true}
		return ma.NewMessage(p.question.Text, ma.REGULAR, msg, keyboard, nil)
	}

	return ma.NewMessage(p.question.Text, ma.REGULAR, msg, nil, nil)
}

func (p *Prompt) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {

	if msg.Type == ma.CALLBACK {
		return nil, EmptyStep
	}

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		return nil, PreviousStep
	}

	p.inProgress = false
	p.State.RawInput[p.question.Field] = msg.Text
	return nil, p.nextStep
}

type Test struct {
	StepBase
}

func (t *Test) Request(msg *ma.Message) *ma.Message {
	t.inProgress = true

	row1 := []tgbotapi.KeyboardButton{
		{Text: "WebApp1", WebApp: &tgbotapi.WebAppInfo{Url: "https://bot-dev-domain.com/webapp1.html"}},
		{Text: "WebApp2", WebApp: &tgbotapi.WebAppInfo{Url: "https://bot-dev-domain.com/webapp2.html"}},
		{Text: "Назад"},
	}

	var keyboard [][]tgbotapi.KeyboardButton

	keyboard = append(keyboard, row1)

	numericKeyboard := &tgbotapi.ReplyKeyboardMarkup{
		Keyboard:       keyboard,
		ResizeKeyboard: true,
	}

	return ma.NewMessage("WebApp test step", ma.REGULAR, msg, numericKeyboard, nil)
}

func (t *Test) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	t.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		return nil, PreviousStep
	}

	return nil, EmptyStep
}

func (t *Test) IsCallBackStep() bool {
	return false
}
