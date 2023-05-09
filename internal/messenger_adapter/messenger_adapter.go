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
	TELEGRAM_CALLBACK
)

type UserData struct {
	WaData     types.MessageInfo
	TgData     *tgbotapi.Message
	TgCallback *tgbotapi.CallbackQuery
}

type TgMarkup struct {
	ReplyMarkup  *tgbotapi.ReplyKeyboardMarkup
	InlineMarkup *tgbotapi.InlineKeyboardMarkup
}

type Message struct {
	Text   string
	Type   int
	UserID string
	UserData
	TgMarkup *TgMarkup
}

func NewMessage(text string, msg *Message, replyMarkup *tgbotapi.ReplyKeyboardMarkup, inlineMarkup *tgbotapi.InlineKeyboardMarkup) *Message {
	return &Message{
		Text:     text,
		UserData: msg.UserData,
		UserID:   msg.UserID,
		Type:     msg.Type,
		TgMarkup: &TgMarkup{
			ReplyMarkup:  replyMarkup,
			InlineMarkup: inlineMarkup,
		},
	}
}
