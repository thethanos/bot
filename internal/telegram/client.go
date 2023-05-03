package telegram

import (
	"fmt"
	"multimessenger_bot/internal/config"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramClient struct {
	client *tgbotapi.BotAPI
	cfg    *config.Config
}

func NewTelegramClient(cfg *config.Config) (*TelegramClient, error) {

	client, err := tgbotapi.NewBotAPI(cfg.TgToken)
	if err != nil {
		return nil, err
	}

	//client.Debug = true

	return &TelegramClient{client: client, cfg: cfg}, nil
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

			fmt.Println(event.Message.Text)
		}
	}()

	return nil
}

func (tc *TelegramClient) Disconnect() {
	tc.client.StopReceivingUpdates()
}

func (tc *TelegramClient) SendMessage(message string) {

}
