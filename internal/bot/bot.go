package bot

import (
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/entities"
	"multimessenger_bot/internal/logger"
	ma "multimessenger_bot/internal/messenger_adapter"
	"strings"
	"time"
)

type UserSession struct {
	CurrentStep  Step
	PrevSteps    StepStack
	State        *entities.UserState
	LastActivity time.Time
}

type Bot struct {
	logger       logger.Logger
	clients      map[ma.MessageSource]ma.ClientInterface
	userSessions map[string]*UserSession
	recvMsgChan  chan *ma.Message
	sendMsgChan  chan *ma.Message
	dbAdapter    *db_adapter.DbAdapter
}

func NewBot(logger logger.Logger, clientArray []ma.ClientInterface, dbAdapter *db_adapter.DbAdapter, recvMsgChan chan *ma.Message) (*Bot, error) {

	clients := make(map[ma.MessageSource]ma.ClientInterface)
	for _, client := range clientArray {
		clients[client.GetType()] = client
	}

	userSessions := make(map[string]*UserSession)
	sendMsgChan := make(chan *ma.Message)

	bot := &Bot{
		logger:       logger,
		clients:      clients,
		userSessions: userSessions,
		recvMsgChan:  recvMsgChan,
		dbAdapter:    dbAdapter,
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
			if _, exists := b.userSessions[msg.UserID]; !exists || strings.ToLower(msg.Text) == "/start" {
				state := &entities.UserState{RawInput: make(map[string]string)}
				b.userSessions[msg.UserID] = &UserSession{
					State:       state,
					CurrentStep: b.createStep(MainMenuStep, state),
					PrevSteps:   StepStack{},
				}
			}
			b.userSessions[msg.UserID].LastActivity = time.Now()
			b.processUserSession(msg)
		}
	}()

	go func() {
		for msg := range b.sendMsgChan {
			if err := b.clients[msg.Source].SendMessage(msg); err != nil {
				b.logger.Error("bot::Run::SendMessage", err)
			}
		}
	}()

	go func() {
		for {
			time.Sleep(time.Hour)
			for id, user := range b.userSessions {
				if time.Now().Sub(user.LastActivity) >= (time.Hour * 24) {
					b.logger.Infof("User session %s has been deleted due to inactivity for the last 24 hours", id)
					delete(b.userSessions, id)
				}
			}
		}
	}()
}

func (b *Bot) Shutdown() {
	for _, client := range b.clients {
		client.Disconnect()
	}
}

func (b *Bot) createStep(step StepType, state *entities.UserState) Step {
	switch step {
	case MainMenuStep:
		return &MainMenu{
			StepBase: StepBase{logger: b.logger, state: state},
		}
	case CitySelectionStep:
		return &CitySelection{
			StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter},
			mode:     &BaseCitySelectionMode{dbAdapter: b.dbAdapter},
		}
	case MainMenuCitySelectionStep:
		return &CitySelection{
			StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter},
			mode: &MainMenuCitySelectionMode{
				BaseCitySelectionMode{
					dbAdapter: b.dbAdapter,
				},
			},
		}
	case ServiceCategorySelectionStep:
		return &ServiceCategorySelection{
			StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter},
			mode: &BaseServiceCategoryMode{
				dbAdapter: b.dbAdapter,
			},
		}
	case MainMenuServiceCategorySelectionStep:
		return &ServiceCategorySelection{
			StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter},
			mode: &MainMenuServiceCategoryMode{
				BaseServiceCategoryMode: BaseServiceCategoryMode{
					dbAdapter: b.dbAdapter,
				},
			},
		}
	case MainMenuServiceSelectionStep:
		return &ServiceSelection{
			StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter},
			mode:     &MainMenuServiceSelectionMode{BaseServiceSelectionMode{dbAdapter: b.dbAdapter}},
		}
	case ServiceSelectionStep:
		return &ServiceSelection{
			StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter},
			mode:     &BaseServiceSelectionMode{dbAdapter: b.dbAdapter},
		}
	case MasterSelectionStep:
		return &MasterSelection{
			StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter},
		}
	case FindModelStep:
		return &FindModel{
			StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter},
		}
	case CollaborationStep:
		return &Collaboration{
			StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter},
		}
	case EmptyStep:
		return nil
	case AdminStep:
		return &Admin{StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter}}
	case AdminServiceCategorySelectionStep:
		return &ServiceCategorySelection{
			StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter},
			mode: &AdminServiceCategoryMode{
				BaseServiceCategoryMode: BaseServiceCategoryMode{dbAdapter: b.dbAdapter},
			},
		}
	case AddServiceCategoryStep:
		return &AddServiceCategory{StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter}}
	case AddServiceStep:
		return &AddService{StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter}}
	case AddCityStep:
		return &AddCity{StepBase: StepBase{logger: b.logger, state: state, dbAdapter: b.dbAdapter}}
	default:
		return &MainMenu{StepBase: StepBase{logger: b.logger, state: state}}
	}
}

func (b *Bot) send(msg *ma.Message) bool {
	if msg == nil {
		return false
	}
	b.sendMsgChan <- msg
	b.logger.Infof("sending a message: %s", msg.Text)
	return true
}

func (b *Bot) processMessage(msg *ma.Message) {
	curStep := b.userSessions[msg.UserID].CurrentStep
	state := b.userSessions[msg.UserID].State
	if !curStep.IsInProgress() {
		b.send(curStep.Request(msg))
	} else {
		res, next := curStep.ProcessResponse(msg)
		b.send(res)

		switch step := b.createStep(next, state); next {
		case PreviousStep:
			var prevStep Step
			if b.userSessions[msg.UserID].PrevSteps.Empty() {
				prevStep = b.createStep(MainMenuStep, state)
			} else {
				prevStep = b.userSessions[msg.UserID].PrevSteps.Top()
				b.userSessions[msg.UserID].PrevSteps.Pop()
			}
			prevStep.Reset()

			b.send(prevStep.Request(msg))
			b.userSessions[msg.UserID].CurrentStep = prevStep
		case EmptyStep:
		case MainMenuStep:
			b.send(step.Request(msg))
			b.userSessions[msg.UserID].CurrentStep = step
			b.userSessions[msg.UserID].PrevSteps.Clear()
		default:
			b.send(step.Request(msg))
			b.userSessions[msg.UserID].CurrentStep = step
			b.userSessions[msg.UserID].PrevSteps.Push(curStep)
		}
	}
}

func (b *Bot) processUserSession(msg *ma.Message) {
	b.processMessage(msg)
}
