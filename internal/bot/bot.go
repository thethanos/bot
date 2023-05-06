package bot

import (
	ma "multimessenger_bot/internal/messenger_adapter"
)

type Bot struct {
	clients      map[int]ma.ClientInterface
	recvMsgChan  chan *ma.Message
	userSessions map[string]*UserSession
	sendMsgChan  chan *ma.Message
}

func NewBot(clientArray []ma.ClientInterface, recvMsgChan chan *ma.Message) (*Bot, error) {

	clients := make(map[int]ma.ClientInterface)
	for _, client := range clientArray {
		clients[client.GetType()] = client
	}

	userSessions := make(map[string]*UserSession)
	sendMsgChan := make(chan *ma.Message)

	return &Bot{clients: clients, recvMsgChan: recvMsgChan, userSessions: userSessions, sendMsgChan: sendMsgChan}, nil
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
		return &CitySelection{StepBase{State: state}}
	case ServiceSelectionStep:
		return &ServiceSelection{StepBase{State: state}}
	case MasterSelectionStep:
		return &MasterSelection{StepBase{State: state}}
	case FinalStep:
		return &Final{StepBase{State: state}}
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
