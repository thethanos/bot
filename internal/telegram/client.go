package telegram

import (
	"fmt"
	"multimessenger_bot/internal/config"
	ma "multimessenger_bot/internal/messenger_adapter"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramClient struct {
	client      *tgbotapi.BotAPI
	cfg         *config.Config
	recvMsgChan chan *ma.Message
}

func NewTelegramClient(cfg *config.Config, recvMsgChan chan *ma.Message) (*TelegramClient, error) {

	client, err := tgbotapi.NewBotAPI(cfg.TgToken)
	if err != nil {
		return nil, err
	}
	return &TelegramClient{client: client, cfg: cfg, recvMsgChan: recvMsgChan}, nil
}

func (tc *TelegramClient) Connect() error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30

	events := tc.client.GetUpdatesChan(updateConfig)
	go func() {
		for event := range events {
			if event.Message != nil {
				msg := &ma.Message{
					Text:     event.Message.Text,
					Type:     ma.TELEGRAM,
					UserID:   fmt.Sprintf("tg%d", event.Message.From.ID),
					UserData: ma.UserData{TgData: event.Message},
				}
				tc.recvMsgChan <- msg
			} else if event.CallbackQuery != nil {
				msg := &ma.Message{
					Type:     ma.TELEGRAM_CALLBACK,
					UserID:   fmt.Sprintf("tg%d", event.CallbackQuery.From.ID),
					UserData: ma.UserData{TgCallback: event.CallbackQuery},
				}
				tc.recvMsgChan <- msg
			}
		}
	}()

	return nil
}

func (tc *TelegramClient) Disconnect() {
	tc.client.StopReceivingUpdates()
}

func (tc *TelegramClient) SendMessage(msg *ma.Message) error {
	if msg == nil {
		return nil
	}

	send := tgbotapi.NewMessage(msg.TgData.From.ID, msg.Text)
	if msg.TgMarkup != nil {
		if msg.TgMarkup.ReplyMarkup != nil {
			send.ReplyMarkup = msg.TgMarkup.ReplyMarkup
		} else {
			send.ReplyMarkup = msg.TgMarkup.InlineMarkup
		}
	} else {
		send.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
	}
	_, err := tc.client.Send(send)
	return err
}

func (tc *TelegramClient) GetType() int {
	return ma.TELEGRAM
}
