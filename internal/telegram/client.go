package telegram

import (
	"multimessenger_bot/internal/config"
	ma "multimessenger_bot/internal/messenger_adapter"
	"net/http"

	handler "multimessenger_bot/internal/telegram/event_handler"

	"github.com/PaulSonOfLars/gotgbot/v2"
	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type TelegramClient struct {
	client      *tgbotapi.Bot
	cfg         *config.Config
	recvMsgChan chan *ma.Message
}

func NewTelegramClient(cfg *config.Config, recvMsgChan chan *ma.Message) (*TelegramClient, error) {

	client, err := tgbotapi.NewBot(cfg.TgToken, &gotgbot.BotOpts{
		Client: http.Client{},
		DefaultRequestOpts: &gotgbot.RequestOpts{
			Timeout: gotgbot.DefaultTimeout,
			APIURL:  gotgbot.DefaultAPIURL,
		},
	})
	if err != nil {
		return nil, err
	}
	return &TelegramClient{client: client, cfg: cfg, recvMsgChan: recvMsgChan}, nil
}

func (tc *TelegramClient) Connect() error {

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{})
	updates := ext.NewUpdater(&ext.UpdaterOpts{Dispatcher: dispatcher})

	updates.StartPolling(tc.client, &ext.PollingOpts{
		DropPendingUpdates: true,
	})

	handler := &handler.Handler{RecvMsgChan: tc.recvMsgChan}
	dispatcher.AddHandler(handler)
	return nil
}

func (tc *TelegramClient) Disconnect() {
}

func (tc *TelegramClient) SendMessage(msg *ma.Message) error {
	if msg == nil {
		return nil
	}

	opts := &tgbotapi.SendMessageOpts{ReplyMarkup: msg.GetTgMarkup()}
	if msg.Type == ma.CALLBACK {
		opts.ParseMode = "HTML"
	}

	_, err := tc.client.SendMessage(msg.GetTgID(), msg.Text, opts)
	return err
}

func (tc *TelegramClient) GetType() ma.MessageSource {
	return ma.TELEGRAM
}
