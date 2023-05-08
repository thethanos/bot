package bot

import (
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/entities"
	ma "multimessenger_bot/internal/messenger_adapter"
	"time"
)

type UserSession struct {
	CurrentStep Step
	State       entities.UserState
}

type Bot struct {
	clients      map[int]ma.ClientInterface
	userSessions map[string]*UserSession
	recvMsgChan  chan *ma.Message
	sendMsgChan  chan *ma.Message
	dbAdapter    *db_adapter.DbAdapter
}

func NewBot(clientArray []ma.ClientInterface, dbAdpter *db_adapter.DbAdapter, recvMsgChan chan *ma.Message) (*Bot, error) {

	clients := make(map[int]ma.ClientInterface)
	for _, client := range clientArray {
		clients[client.GetType()] = client
	}

	userSessions := make(map[string]*UserSession)
	sendMsgChan := make(chan *ma.Message)

	bot := &Bot{
		clients:      clients,
		userSessions: userSessions,
		recvMsgChan:  recvMsgChan,
		dbAdapter:    dbAdpter,
		sendMsgChan:  sendMsgChan,
	}

	return bot, nil
}

func (b *Bot) Run() {

	for _, client := range b.clients {
		client.Connect()
	}

	go func() {
		for msg := range b.recvMsgChan {
			if _, exists := b.userSessions[msg.UserID]; !exists {
				b.userSessions[msg.UserID] = &UserSession{
					State:       entities.UserState{Cursor: 0, RawInput: make(map[string]string)},
					CurrentStep: b.createStep(MainMenuStep, nil),
				}
			}

			b.processUserSession(msg)
		}
	}()

	go func() {
		for msg := range b.sendMsgChan {
			b.clients[msg.Type].SendMessage(msg)
		}
	}()
}

func (b *Bot) Shutdown() {
	for _, client := range b.clients {
		client.Disconnect()
	}
}

func (b *Bot) createStep(step int, state *entities.UserState) Step {
	switch step {
	case MainMenuStep:
		return &MainMenu{}
	case MainMenuRequestStep:
		return &YesNo{question: Question{Text: "Вурнуться в главное меню?"}, yesStep: MainMenuStep, noStep: EmptyStep}
	case CitySelectionStep:
		return &CitySelection{
			StepBase:     StepBase{State: state, DbAdapter: b.dbAdapter},
			checkService: true,
			filter:       true,
			nextStep:     MasterSelectionStep,
			errStep:      EmptyStep,
		}
	case ServiceSelectionStep:
		return &ServiceSelection{
			StepBase:  StepBase{State: state, DbAdapter: b.dbAdapter},
			checkCity: true,
			filter:    true,
			nextStep:  MasterSelectionStep,
			errStep:   EmptyStep,
		}
	case MasterSelectionStep:
		return &MasterSelection{StepBase: StepBase{State: state, DbAdapter: b.dbAdapter}}
	case FinalStep:
		return &Final{StepBase{State: state, DbAdapter: b.dbAdapter}}
	case MasterStep:
		return &YesNo{
			StepBase: StepBase{State: state, DbAdapter: b.dbAdapter},
			question: Question{Text: "Хотите зарегистрироваться как мастер?"},
			yesStep:  RegistrationStep,
			noStep:   MainMenuStep,
		}
	case RegistrationStep:
		return &Prompt{
			StepBase: StepBase{State: state, DbAdapter: b.dbAdapter},
			question: Question{Text: "Как вас называть?", Field: "name"},
			nextStep: RegistrationStepCity,
			errStep:  RegistrationStep,
		}
	case RegistrationStepCity:
		return &CitySelection{
			StepBase: StepBase{State: state, DbAdapter: b.dbAdapter},
			nextStep: RegistrationStepService,
			errStep:  EmptyStep,
		}
	case RegistrationStepService:
		return &ServiceSelection{
			StepBase: StepBase{State: state, DbAdapter: b.dbAdapter},
			nextStep: RegistrationFinalStep,
			errStep:  EmptyStep,
		}
	case RegistrationFinalStep:
		return &RegistrationFinal{StepBase: StepBase{State: state, DbAdapter: b.dbAdapter}}
	case EmptyStep:
		return nil
	default:
		return &MainMenu{}
	}
}

func (b *Bot) send(msg *ma.Message) bool {
	if msg == nil {
		return false
	}
	b.sendMsgChan <- msg
	return true
}

func (b *Bot) processUserSession(msg *ma.Message) {
	curStep := b.userSessions[msg.UserID].CurrentStep
	state := &b.userSessions[msg.UserID].State
	if !curStep.IsInProgress() {
		b.send(curStep.Request(msg))
	} else {
		res, next := curStep.ProcessResponse(msg)
		b.send(res)

		switch step := b.createStep(next, state); next {
		case EmptyStep:
		case MainMenuRequestStep:
			time.Sleep(1 * time.Second)
			fallthrough
		default:
			b.send(step.Request(msg))
			b.userSessions[msg.UserID].CurrentStep = step
		}
	}
}
