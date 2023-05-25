package messenger_adapter

import (
	"os"

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
	DownloadFile(*Message)
}

type MessageData struct {
	WaData       types.MessageInfo
	TgData       *tgbotapi.Message
	TgMarkup     *tgbotapi.ReplyKeyboardMarkup
	removeMarkup bool
}

type Message struct {
	Text   string
	Image  []byte
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
	if m.Data.removeMarkup {
		return tgbotapi.ReplyKeyboardRemove{RemoveKeyboard: true}
	}
	return m.Data.TgMarkup
}

func NewTextMessage(text string, msg *Message, replyMarkup *tgbotapi.ReplyKeyboardMarkup, removeMarkup bool) *Message {

	if msg.Data == nil {
		panic("Empty data")
	}

	data := &MessageData{
		WaData:       msg.Data.WaData,
		TgData:       msg.Data.TgData,
		TgMarkup:     replyMarkup,
		removeMarkup: removeMarkup,
	}

	return &Message{
		Text:   text,
		Type:   TEXT,
		UserID: msg.UserID,
		Source: msg.Source,
		Data:   data,
	}
}

func NewImageMessage(path, caption string, msg *Message, removeMarkup bool) *Message {

	image, err := os.ReadFile(path)
	if err != nil {
		panic("Failed to load image")
	}

	data := &MessageData{
		WaData:       msg.Data.WaData,
		TgData:       msg.Data.TgData,
		removeMarkup: removeMarkup,
	}

	return &Message{
		Text:   caption,
		Image:  image,
		Type:   IMAGE,
		UserID: msg.UserID,
		Source: msg.Source,
		Data:   data,
	}
}
