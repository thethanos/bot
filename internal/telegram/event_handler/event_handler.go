package telegram

import (
	"fmt"
	ma "multimessenger_bot/internal/messenger_adapter"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Handler struct {
	RecvMsgChan chan *ma.Message
}

func (h *Handler) CheckUpdate(client *tgbotapi.Bot, ctx *ext.Context) bool {
	return true
}

func (h *Handler) HandleUpdate(client *tgbotapi.Bot, ctx *ext.Context) error {
	event := ctx.Update
	if event.Message != nil {
		msg := &ma.Message{
			Text:   event.Message.Text,
			Type:   ma.REGULAR,
			Source: ma.TELEGRAM,
			UserID: fmt.Sprintf("tg%d", event.Message.From.Id),
			Data:   &ma.MessageData{TgData: event.Message},
		}
		h.RecvMsgChan <- msg
	} else if event.CallbackQuery != nil {
		msg := &ma.Message{
			Type:   ma.CALLBACK,
			Source: ma.TELEGRAM,
			UserID: fmt.Sprintf("tg%d", event.CallbackQuery.From.Id),
			Data:   &ma.MessageData{TgCallback: event.CallbackQuery},
		}
		h.RecvMsgChan <- msg
	}
	return nil
}

func (h *Handler) Name() string {
	return "custom handler"
}
