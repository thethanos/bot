package whatsapp

import (
	"fmt"
	ma "multimessenger_bot/internal/messenger_adapter"

	"go.mau.fi/whatsmeow/types/events"
)

type Handler struct {
	RecvMsgChan chan *ma.Message
}

func (h *Handler) EventHandler(event interface{}) {
	switch v := event.(type) {
	case *events.Message:
		userId := fmt.Sprintf("wa%s", v.Info.Chat.User)
		h.RecvMsgChan <- &ma.Message{Text: v.Message.GetConversation(), Type: ma.WHATSAPP, UserID: userId, UserData: ma.UserData{WaData: v.Info}}
	}
}
