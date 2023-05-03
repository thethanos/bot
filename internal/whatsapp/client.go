package whatsapp

import (
	"context"
	"os"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"

	"multimessenger_bot/internal/config"
	handler "multimessenger_bot/internal/whatsapp/event_handler"

	"github.com/mdp/qrterminal"
)

type DeviceManager interface {
	GetFirstDevice() (*store.Device, error)
}

type WhatsAppClient struct {
	client *whatsmeow.Client
	cfg    *config.Config
}

func NewWhatsAppClient(log waLog.Logger, cfg *config.Config, dm DeviceManager) (*WhatsAppClient, error) {

	deviceStore, err := dm.GetFirstDevice()
	if err != nil {
		return nil, err
	}

	client := whatsmeow.NewClient(deviceStore, log)
	client.AddEventHandler(handler.EventHandler)

	return &WhatsAppClient{client: client, cfg: cfg}, nil
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

func (wc *WhatsAppClient) SendMessage(message string) {

}
