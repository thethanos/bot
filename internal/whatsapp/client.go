package whatsapp

import (
	"context"
	"os"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	waLog "go.mau.fi/whatsmeow/util/log"

	"bot/internal/config"
	ma "bot/internal/msgadapter"
	handler "bot/internal/whatsapp/event_handler"

	"github.com/mdp/qrterminal"
)

type DeviceManager interface {
	GetFirstDevice() (*store.Device, error)
}

type WhatsAppClient struct {
	client      *whatsmeow.Client
	cfg         *config.Config
	recvMsgChan chan *ma.Message
}

func NewWhatsAppClient(log waLog.Logger, cfg *config.Config, dm DeviceManager, recvMsgChan chan *ma.Message) (*WhatsAppClient, error) {

	deviceStore, err := dm.GetFirstDevice()
	if err != nil {
		return nil, err
	}

	client := whatsmeow.NewClient(deviceStore, log)
	handler := handler.Handler{RecvMsgChan: recvMsgChan}
	client.AddEventHandler(handler.EventHandler)

	return &WhatsAppClient{client: client, cfg: cfg, recvMsgChan: recvMsgChan}, nil
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

func (wc *WhatsAppClient) SendMessage(msg *ma.Message) error {
	if msg == nil {
		return nil
	}
	toSend := &proto.Message{Conversation: &msg.Text}

	_, err := wc.client.SendMessage(context.Background(), msg.GetWaID(), toSend)
	return err
}

func (wc *WhatsAppClient) GetType() ma.MessageSource {
	return ma.WHATSAPP
}

func (wc *WhatsAppClient) DownloadFile(id string, msg *ma.Message) string {
	return ""
}
