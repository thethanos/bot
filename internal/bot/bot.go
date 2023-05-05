package bot

import (
	"fmt"
	ci "multimessenger_bot/internal/client_interface"
)

type Bot struct {
	clients      map[int]ci.ClientInterface
	msgChan      chan ci.Message
	userSessions map[string]*UserSession
	sendMsgChan  chan ci.Message
}

func NewBot(clientArray []ci.ClientInterface, msgChan chan ci.Message) (*Bot, error) {

	clients := make(map[int]ci.ClientInterface)
	for _, client := range clientArray {
		clients[client.GetType()] = client
	}

	userSessions := make(map[string]*UserSession)
	sendMsgChan := make(chan ci.Message)

	return &Bot{clients: clients, msgChan: msgChan, userSessions: userSessions, sendMsgChan: sendMsgChan}, nil
}

func (b *Bot) Run() {

	for _, client := range b.clients {
		client.Connect()
	}

	go func() {
		for msg := range b.msgChan {
			fmt.Println(msg.Text)

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
		return &MainMenu{State: nil}
	case CitiesStep:
		return &Cities{State: state}
	case ServicesStep:
		return &Services{State: state}
	case MasterStep:
		return &Master{State: state}
	case FinalStep:
		return &Final{State: state}
	case EmptyStep:
		return &Empty{}
	default:
		return &MainMenu{State: nil}
	}
}

func (b *Bot) processUserSession(msg ci.Message) {
	curStep := b.userSessions[msg.UserID].CurrentStep
	state := &b.userSessions[msg.UserID].State
	if !curStep.IsInProgress() {
		b.sendMsgChan <- curStep.Request(msg)
	} else {
		res, next := curStep.ProcessResponse(msg)
		b.sendMsgChan <- res

		step := b.createStep(next, state)

		switch next {
		case MainMenuStep:
			b.sendMsgChan <- step.DefaultRequest(msg)
			b.userSessions[msg.UserID].CurrentStep = step
		case EmptyStep:
		default:
			b.sendMsgChan <- step.Request(msg)
			b.userSessions[msg.UserID].CurrentStep = step
		}
	}
}
