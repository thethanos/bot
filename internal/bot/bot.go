package bot

import (
	"multimessenger_bot/internal/db_adapter"
	ma "multimessenger_bot/internal/messenger_adapter"
)

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
				b.userSessions[msg.UserID] = &UserSession{CurrentStep: b.createStep(MainMenuStep, nil)}
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

func (b *Bot) createStep(step int, state *UserState) Step {
	switch step {
	case MainMenuStep:
		return &MainMenu{}
	case CitySelectionStep:
		return &CitySelection{StepBase: StepBase{State: state, DbAdapter: b.dbAdapter}}
	case ServiceSelectionStep:
		return &ServiceSelection{StepBase: StepBase{State: state, DbAdapter: b.dbAdapter}}
	case MasterSelectionStep:
		return &MasterSelection{StepBase: StepBase{State: state, DbAdapter: b.dbAdapter}}
	case FinalStep:
		return &Final{StepBase{State: state, DbAdapter: b.dbAdapter}}
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
		default:
			b.send(step.Request(msg))
			b.userSessions[msg.UserID].CurrentStep = step
		}
	}
}
