package whatsapp

import (
	"fmt"
	ma "multimessenger_bot/internal/msgadapter"

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
