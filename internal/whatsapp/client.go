package whatsapp

import (
	"context"
	"fmt"
	"os"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"

	"multimessenger_bot/internal/client_interface"
	ci "multimessenger_bot/internal/client_interface"
	"multimessenger_bot/internal/config"
	handler "multimessenger_bot/internal/whatsapp/event_handler"

	"github.com/mdp/qrterminal"
)

type DeviceManager interface {
	GetFirstDevice() (*store.Device, error)
}

type WhatsAppClient struct {
	client  *whatsmeow.Client
	cfg     *config.Config
	msgChan chan client_interface.Message
}

func NewWhatsAppClient(log waLog.Logger, cfg *config.Config, dm DeviceManager, msgChan chan client_interface.Message) (*WhatsAppClient, error) {

	deviceStore, err := dm.GetFirstDevice()
	if err != nil {
		return nil, err
	}

	client := whatsmeow.NewClient(deviceStore, log)
	handler := handler.Handler{MsgChan: msgChan}
	client.AddEventHandler(handler.EventHandler)

	return &WhatsAppClient{client: client, cfg: cfg, msgChan: msgChan}, nil
}

func (wc *WhatsAppClient) Connect() error {
	if wc.client.Store.ID == nil {
		qrChan, _ := wc.client.GetQRChannel(context.Background())
		if err := wc.client.Connect(); err != nil {
			return err
		}

		for event := range qrChan {
			if event.Event == "code" {
				qrterminal.GenerateHalfBlock(event.Code, qrterminal.L, os.Stdout)
			}
		}
	} else {
		if err := wc.client.Connect(); err != nil {
			return err
		}
	}
	return nil
}

func (wc *WhatsAppClient) Disconnect() {
	wc.client.Disconnect()
}

func (wc *WhatsAppClient) SendMessage(msg ci.Message) {
	if len(msg.Text) == 0 {
		return
	}
	toSend := &proto.Message{Conversation: &msg.Text}

	resp, err := wc.client.SendMessage(context.Background(), msg.WaData.Chat, toSend)
	fmt.Println(resp, err)
}

func (wc *WhatsAppClient) GetType() int {
	return ci.WHATSAPP
}
