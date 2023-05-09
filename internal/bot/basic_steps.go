package bot

import (
	"fmt"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/entities"
	ma "multimessenger_bot/internal/messenger_adapter"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	MainMenuStep = iota
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
	ProcessResponse(*ma.Message) (*ma.Message, int)
	Request(*ma.Message) *ma.Message
	IsInProgress() bool
	Reset()
}

type StepBase struct {
	inProgress bool
	State      *entities.UserState
	DbAdapter  *db_adapter.DbAdapter
}

func (s *StepBase) IsInProgress() bool {
	return s.inProgress
}

func (s *StepBase) Reset() {
}

type YesNo struct {
	StepBase
	question Question
	yesStep  int
	noStep   int
}

func (y *YesNo) Request(msg *ma.Message) *ma.Message {
	y.inProgress = true
	if msg.Type == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 2)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Да"}}
		rows[1] = []tgbotapi.KeyboardButton{{Text: "Нет"}}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return &ma.Message{Text: y.question.Text, UserData: msg.UserData, Type: msg.Type, TgMarkup: keyboard}
	}
	return &ma.Message{Text: fmt.Sprintf("%s\n1. Да\n2. Нет", y.question.Text), UserData: msg.UserData, Type: msg.Type}
}

func (y *YesNo) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
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
	nextStep int
	errStep  int
}

func (p *Prompt) Request(msg *ma.Message) *ma.Message {
	p.inProgress = true
	return &ma.Message{Text: p.question.Text, UserData: msg.UserData, Type: msg.Type}
}

func (p *Prompt) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	p.inProgress = false
	p.State.RawInput[p.question.Field] = msg.Text
	return nil, p.nextStep
}
