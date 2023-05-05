package whatsapp

import (
	ci "multimessenger_bot/internal/client_interface"

	"go.mau.fi/whatsmeow/types/events"
)

type Handler struct {
	MsgChan chan ci.Message
}

func (h *Handler) EventHandler(event interface{}) {
	switch v := event.(type) {
	case *events.Message:
		h.MsgChan <- ci.Message{Text: v.Message.GetConversation(), Type: ci.WHATSAPP, WaData: v.Info}
	}
}
