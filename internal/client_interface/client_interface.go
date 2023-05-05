package client_interface

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mau.fi/whatsmeow/types"
)

type ClientInterface interface {
	Connect() error
	Disconnect()
	SendMessage(Message)
	GetType() int
}

const (
	WHATSAPP = iota
	TELEGRAM
)

type Message struct {
	Text   string
	Type   int
	UserID string
	WaData types.MessageInfo
	TgData tgbotapi.Message
}
