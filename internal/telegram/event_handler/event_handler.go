package telegram

import (
	"bot/internal/logger"
	ma "bot/internal/msgadapter"
	"fmt"

	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type Handler struct {
	logger      logger.Logger
	recvMsgChan chan *ma.Message
}

func NewHandler(logger logger.Logger, recvMsgChan chan *ma.Message) *Handler {
	return &Handler{
		logger:      logger,
		recvMsgChan: recvMsgChan,
	}
}

func (h *Handler) CheckUpdate(client *tgbotapi.Bot, ctx *ext.Context) bool {
	return true
}

func (h *Handler) HandleUpdate(client *tgbotapi.Bot, ctx *ext.Context) error {
	event := ctx.Update
	if event.Message != nil {
		msg := &ma.Message{
			Text:   event.Message.Text,
			Source: ma.TELEGRAM,
			UserID: fmt.Sprintf("tg%d", event.Message.From.Id),
			Data:   &ma.MessageData{TgData: event.Message},
		}
		h.recvMsgChan <- msg
	}
	return nil
}

func (h *Handler) Name() string {
	return "custom handler"
}
