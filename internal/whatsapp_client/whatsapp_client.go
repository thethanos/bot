package whatsapp_client

import (
	"context"
	"os"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"

	handler "whatsapp_bot/internal/whatsapp_client/event_handler"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mdp/qrterminal"
)

type DeviceManager interface {
	GetFirstDevice() (*store.Device, error)
}

type WhatsAppClient struct {
	client *whatsmeow.Client
}

func NewWhatsAppClient(log waLog.Logger, dm DeviceManager) *WhatsAppClient {

	deviceStore, err := dm.GetFirstDevice()
	if err != nil {
		panic(err)
	}

	client := whatsmeow.NewClient(deviceStore, log)
	client.AddEventHandler(handler.EventHandler)

	return &WhatsAppClient{client: client}
}

func (wc *WhatsAppClient) Connect() {
	if wc.client.Store.ID == nil {
		qrChan, _ := wc.client.GetQRChannel(context.Background())
		if err := wc.client.Connect(); err != nil {
			panic(err)
		}

		for event := range qrChan {
			if event.Event == "code" {
				qrterminal.GenerateHalfBlock(event.Code, qrterminal.L, os.Stdout)
			}
		}
	} else {
		if err := wc.client.Connect(); err != nil {
			panic(err)
		}
	}
}

func (wc *WhatsAppClient) Disconnect() {
	wc.client.Disconnect()
}
