package messenger_adapter

import (
	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
	"go.mau.fi/whatsmeow/types"
)

type MessageType uint

const (
	TEXT MessageType = iota
	IMAGE
)

type MessageSource uint

const (
	WHATSAPP MessageSource = iota
	TELEGRAM
)

type ClientInterface interface {
	Connect() error
	Disconnect()
	SendMessage(*Message) error
	GetType() MessageSource
}

type MessageData struct {
	WaData   types.MessageInfo
	TgData   *tgbotapi.Message
	TgMarkup *tgbotapi.ReplyKeyboardMarkup
}

type Message struct {
	Text   string
	Type   MessageType
	Source MessageSource
	UserID string
	Data   *MessageData
}

func (m *Message) GetTgID() int64 {
	if m.Data.TgData != nil {
		return m.Data.TgData.From.Id
	}
	return 0
}

func (m *Message) GetWaID() types.JID {
	return m.Data.WaData.Chat
}

func (m *Message) GetTgMarkup() tgbotapi.ReplyMarkup {
	if markup := m.Data.TgMarkup; markup != nil {
		return markup
	}
	return tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
}

func NewTextMessage(text string, msg *Message, replyMarkup *tgbotapi.ReplyKeyboardMarkup) *Message {

	if msg.Source == TELEGRAM && msg.Data.TgData == nil {
		panic("Empty data")
	}

	data := &MessageData{
		WaData:   msg.Data.WaData,
		TgData:   msg.Data.TgData,
		TgMarkup: replyMarkup,
	}

	return &Message{
		Text:   text,
		Type:   TEXT,
		UserID: msg.UserID,
		Source: msg.Source,
		Data:   data,
	}
}

func NewImageMessage() *Message {
	return nil
}
