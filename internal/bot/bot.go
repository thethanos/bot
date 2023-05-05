package bot

import (
	"fmt"
	ci "multimessenger_bot/internal/client_interface"
)

type Bot struct {
	clients map[int]ci.ClientInterface
	msgChan chan ci.Message
}

func NewBot(clientArray []ci.ClientInterface, msgChan chan ci.Message) (*Bot, error) {

	clients := make(map[int]ci.ClientInterface)
	for _, client := range clientArray {
		clients[client.GetType()] = client
	}

	return &Bot{clients: clients, msgChan: msgChan}, nil
}

func (b *Bot) Run() {

	for _, client := range b.clients {
		client.Connect()
	}

	go func() {
		for msg := range b.msgChan {
			fmt.Println(msg)
			response := b.processMessage(msg)
			b.clients[msg.Type].SendMessage(response)
		}
	}()
}

func (b *Bot) Shutdown() {
	for _, client := range b.clients {
		client.Disconnect()
	}
}

func (b *Bot) processMessage(msg ci.Message) ci.Message {

	switch msg.Text {
	case "услуги":
		return ci.Message{Text: "выбор услуги: 1) 2) 3)", WaData: msg.WaData, TgData: msg.TgData}
	case "1", "2", "3":
		return ci.Message{Text: "выберите город: город1, город2", WaData: msg.WaData, TgData: msg.TgData}
	case "город1":
		return ci.Message{Text: "finish", WaData: msg.WaData, TgData: msg.TgData}
	default:
		return ci.Message{Text: "error", WaData: msg.WaData, TgData: msg.TgData}
	}
}
