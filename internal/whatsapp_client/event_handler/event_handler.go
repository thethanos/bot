package whatsapp_client

import (
	"fmt"

	"go.mau.fi/whatsmeow/types/events"
)

func EventHandler(event interface{}) {
	switch v := event.(type) {
	case *events.Message:
		fmt.Println("Received a message!", v.Message.GetConversation())
	}
}
