package telegram

import (
	"fmt"
	ci "multimessenger_bot/internal/client_interface"
	"multimessenger_bot/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramClient struct {
	client  *tgbotapi.BotAPI
	cfg     *config.Config
	msgChan chan ci.Message
}

func NewTelegramClient(cfg *config.Config, msgChan chan ci.Message) (*TelegramClient, error) {

	client, err := tgbotapi.NewBotAPI(cfg.TgToken)
	if err != nil {
		return nil, err
	}
	return &TelegramClient{client: client, cfg: cfg, msgChan: msgChan}, nil
}

func (tc *TelegramClient) Connect() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	events := tc.client.GetUpdatesChan(updateConfig)
	go func() {
		for event := range events {
			if event.Message == nil {
				continue
			}
			userId := fmt.Sprintf("tg%d", event.Message.From.ID)
			tc.msgChan <- ci.Message{Text: event.Message.Text, Type: ci.TELEGRAM, UserID: userId, TgData: *event.Message}
		}
	}()

	return nil
}

func (tc *TelegramClient) Disconnect() {
	tc.client.StopReceivingUpdates()
}

func (tc *TelegramClient) SendMessage(msg ci.Message) {
	if len(msg.Text) == 0 {
		return
	}
	tc.client.Send(tgbotapi.NewMessage(msg.TgData.From.ID, msg.Text))
}

func (tc *TelegramClient) GetType() int {
	return ci.TELEGRAM
}
