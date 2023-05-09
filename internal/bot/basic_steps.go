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
		return ma.NewMessage(y.question.Text, msg, keyboard, nil)
	}
	return ma.NewMessage(fmt.Sprintf("%s\n1. Да\n2. Нет", y.question.Text), msg, nil, nil)
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
	return ma.NewMessage(p.question.Text, msg, nil, nil)
}

func (p *Prompt) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	p.inProgress = false
	p.State.RawInput[p.question.Field] = msg.Text
	return nil, p.nextStep
}

type Test struct {
	StepBase
}

func (t *Test) Request(msg *ma.Message) *ma.Message {
	t.inProgress = true
	row1 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	)

	row2 := tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	)

	var keyboard [][]tgbotapi.InlineKeyboardButton

	keyboard = append(keyboard, row1)
	keyboard = append(keyboard, row2)

	numericKeyboard := &tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: keyboard,
	}

	return ma.NewMessage("text", msg, nil, numericKeyboard)
}

func (t *Test) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	t.inProgress = false
	return nil, EmptyStep
}
