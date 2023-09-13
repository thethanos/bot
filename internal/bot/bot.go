package bot

import (
	"bot/internal/dbadapter"
	"bot/internal/entities"
	"bot/internal/logger"
	ma "bot/internal/msgadapter"
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
	DBAdapter    *dbadapter.DBAdapter
}

func NewBot(logger logger.Logger, clientArray []ma.ClientInterface, DBAdapter *dbadapter.DBAdapter, recvMsgChan chan *ma.Message) (*Bot, error) {

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
		DBAdapter:    DBAdapter,
		sendMsgChan:  sendMsgChan,
	}

	return bot, nil
}

func (b *Bot) Run() {

	for _, client := range b.clients {
		if err := client.Connect(); err != nil {
			b.logger.Error("bot::Run::Connect", err)
		}
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
				if time.Since(user.LastActivity) >= (time.Hour * 24) {
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
			StepBase: StepBase{logger: b.logger, state: state, DBAdapter: b.DBAdapter},
			mode:     &BaseCitySelectionMode{dbAdapter: b.DBAdapter},
		}
	case MainMenuCitySelectionStep:
		return &CitySelection{
			StepBase: StepBase{logger: b.logger, state: state, DBAdapter: b.DBAdapter},
			mode: &MainMenuCitySelectionMode{
				BaseCitySelectionMode{
					dbAdapter: b.DBAdapter,
				},
			},
		}
	case ServiceCategorySelectionStep:
		return &ServiceCategorySelection{
			StepBase: StepBase{logger: b.logger, state: state, DBAdapter: b.DBAdapter},
			mode: &BaseServiceCategoryMode{
				dbAdapter: b.DBAdapter,
			},
		}
	case MainMenuServiceCategorySelectionStep:
		return &ServiceCategorySelection{
			StepBase: StepBase{logger: b.logger, state: state, DBAdapter: b.DBAdapter},
			mode: &MainMenuServiceCategoryMode{
				BaseServiceCategoryMode: BaseServiceCategoryMode{
					dbAdapter: b.DBAdapter,
				},
			},
		}
	case ServiceSelectionStep:
		return &ServiceSelection{
			StepBase: StepBase{logger: b.logger, state: state, DBAdapter: b.DBAdapter},
			mode:     &BaseServiceSelectionMode{dbAdapter: b.DBAdapter},
		}
	case MainMenuServiceSelectionStep:
		return &ServiceSelection{
			StepBase: StepBase{logger: b.logger, state: state, DBAdapter: b.DBAdapter},
			mode:     &MainMenuServiceSelectionMode{BaseServiceSelectionMode{dbAdapter: b.DBAdapter}},
		}
	case MasterSelectionStep:
		return &MasterSelection{
			StepBase: StepBase{logger: b.logger, state: state, DBAdapter: b.DBAdapter},
		}
	case FindModelStep:
		return &FindModel{
			StepBase: StepBase{logger: b.logger, state: state, DBAdapter: b.DBAdapter},
		}
	case CollaborationStep:
		return &Collaboration{
			StepBase: StepBase{logger: b.logger, state: state, DBAdapter: b.DBAdapter},
		}
	case EmptyStep:
		return nil
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
