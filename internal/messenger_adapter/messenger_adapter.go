package messenger_adapter

import (
	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
	"go.mau.fi/whatsmeow/types"
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

type MessageType uint

const (
	REGULAR MessageType = iota
	CALLBACK
	WEBAPP
)

type TgMarkup struct {
	ReplyMarkup  *tgbotapi.ReplyKeyboardMarkup
	InlineMarkup *tgbotapi.InlineKeyboardMarkup
}

type MessageData struct {
	WaData     types.MessageInfo
	TgData     *tgbotapi.Message
	TgCallback *tgbotapi.CallbackQuery
	TgMarkup   *TgMarkup
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
	if m.Data.TgCallback != nil {
		return m.Data.TgCallback.From.Id
	}

	return 0
}

func (m *Message) GetWaID() types.JID {
	return m.Data.WaData.Chat
}

func (m *Message) GetTgMarkup() tgbotapi.ReplyMarkup {
	if markup := m.Data.TgMarkup; markup != nil {
		if reply := markup.ReplyMarkup; reply != nil {
			return reply
		}
		if inline := markup.InlineMarkup; inline != nil {
			return inline
		}
	}

	if m.Type == CALLBACK {
		return nil
	} else {
		return tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
	}
}

func NewMessage(text string, msgType MessageType, msg *Message, replyMarkup *tgbotapi.ReplyKeyboardMarkup, inlineMarkup *tgbotapi.InlineKeyboardMarkup) *Message {

	if msg.Type == CALLBACK && msg.Data.TgCallback == nil {
		panic("Empty data")
	}

	if msg.Type == REGULAR && msg.Source == TELEGRAM && msg.Data.TgData == nil {
		panic("Empty data")
	}

	data := &MessageData{
		WaData:     msg.Data.WaData,
		TgData:     msg.Data.TgData,
		TgCallback: msg.Data.TgCallback,
		TgMarkup: &TgMarkup{
			ReplyMarkup:  replyMarkup,
			InlineMarkup: inlineMarkup,
		},
	}

	return &Message{
		Text:   text,
		UserID: msg.UserID,
		Type:   msgType,
		Source: msg.Source,
		Data:   data,
	}
}
