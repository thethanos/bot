package messenger_adapter

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.mau.fi/whatsmeow/types"
)

type ClientInterface interface {
	Connect() error
	Disconnect()
	SendMessage(*Message) error
	GetType() int
}

const (
	WHATSAPP = iota
	TELEGRAM
)

type UserData struct {
	WaData types.MessageInfo
	TgData tgbotapi.Message
}

type Message struct {
	Text   string
	Type   int
	UserID string
	UserData
}
