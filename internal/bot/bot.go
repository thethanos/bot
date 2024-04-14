package bot

import (
	"bot/internal/config"
	"bot/internal/dbadapter"
	"bot/internal/entities"
	"bot/internal/logger"
	client "bot/internal/messenger_client"
	ma "bot/internal/msgadapter"
	"context"
	"strings"
	"sync"
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
	cfg          *config.Config
	clients      map[ma.MessageSource]client.ClientInterface
	userSessions map[string]*UserSession
	recvMsgChan  chan *ma.Message
	sendMsgChan  chan *ma.Message
	DBAdapter    dbadapter.DBInterface
}

func NewBot(logger logger.Logger, cfg *config.Config, clientArray []client.ClientInterface, DBAdapter dbadapter.DBInterface, recvMsgChan chan *ma.Message) (*Bot, error) {

	clients := make(map[ma.MessageSource]client.ClientInterface)
	for _, client := range clientArray {
		clients[client.GetType()] = client
	}

	userSessions := make(map[string]*UserSession)
	sendMsgChan := make(chan *ma.Message, cfg.SendBufSize)

	bot := &Bot{
		logger:       logger,
		cfg:          cfg,
		clients:      clients,
		userSessions: userSessions,
		recvMsgChan:  recvMsgChan,
		DBAdapter:    DBAdapter,
		sendMsgChan:  sendMsgChan,
	}

	return bot, nil
}

func (b *Bot) Run(ctx context.Context, wg *sync.WaitGroup) {

	for _, client := range b.clients {
		if err := client.Connect(); err != nil {
			b.logger.Error("bot::Run::Connect", err)
		}
	}

	go b.processMessages(ctx, wg)
	go b.processSend(ctx, wg)
	go b.cleanUpUserSessions()
}

func (b *Bot) Shutdown() {
	for _, client := range b.clients {
		client.Disconnect()
	}
}

func (b *Bot) processMessages(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-b.recvMsgChan:
			if _, exists := b.userSessions[msg.UserID]; !exists || strings.ToLower(msg.Text) == "/start" {
				state := &entities.UserState{RawInput: make(map[string]string)}
				b.userSessions[msg.UserID] = &UserSession{
					State:       state,
					CurrentStep: b.createStep(MainMenuStep, state),
					PrevSteps:   StepStack{},
				}
				b.send(b.userSessions[msg.UserID].CurrentStep.Request(msg))
			} else {
				b.processUserSession(msg)
			}
			b.userSessions[msg.UserID].LastActivity = time.Now()
		}
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
			StepBase:   StepBase{logger: b.logger, state: state, DBAdapter: b.DBAdapter},
			GalleryURL: b.cfg.GalleryURL,
		}
	case FindModelStep:
		return &FindModel{
			StepBase:  StepBase{logger: b.logger, state: state, DBAdapter: b.DBAdapter},
			ModelsURL: b.cfg.ModelsURL,
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

func (b *Bot) processSend(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-b.sendMsgChan:
			if err := b.clients[msg.Source].SendMessage(msg); err != nil {
				b.logger.Error("bot::Run::SendMessage", err)
			}
		}
	}
}

func (b *Bot) processUserSession(msg *ma.Message) {
	curStep := b.userSessions[msg.UserID].CurrentStep
	state := b.userSessions[msg.UserID].State

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
			prevStep.Reset()
		}

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

func (b *Bot) cleanUpUserSessions() {
	for {
		time.Sleep(time.Hour)
		for id, user := range b.userSessions {
			if time.Since(user.LastActivity) >= (time.Hour * 24) {
				b.logger.Infof("User session %s has been deleted due to inactivity for the last 24 hours", id)
				delete(b.userSessions, id)
			}
		}
	}
}
