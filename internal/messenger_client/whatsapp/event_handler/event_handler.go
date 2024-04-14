package whatsapp

import (
	ma "bot/internal/msgadapter"
	"fmt"

	"go.mau.fi/whatsmeow/types/events"
)

type Handler struct {
	RecvMsgChan chan *ma.Message
}

func (h *Handler) EventHandler(event interface{}) {
	switch v := event.(type) {
	case *events.Message:
		userId := fmt.Sprintf("wa%s", v.Info.Chat.User)
		msg := &ma.Message{
			Text:   v.Message.GetConversation(),
			Source: ma.WHATSAPP,
			UserID: userId,
			Data: &ma.MessageData{
				WaData: v.Info,
			},
		}
		h.RecvMsgChan <- msg
	}
}
