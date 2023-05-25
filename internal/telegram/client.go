package telegram

import (
	"fmt"
	"io"
	"multimessenger_bot/internal/config"
	ma "multimessenger_bot/internal/messenger_adapter"
	"net/http"
	"os"

	handler "multimessenger_bot/internal/telegram/event_handler"

	"github.com/PaulSonOfLars/gotgbot/v2"
	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
	"go.uber.org/zap"
)

type TelegramClient struct {
	logger      *zap.SugaredLogger
	cfg         *config.Config
	recvMsgChan chan *ma.Message
	client      *tgbotapi.Bot
}

func NewTelegramClient(logger *zap.SugaredLogger, cfg *config.Config, recvMsgChan chan *ma.Message) (*TelegramClient, error) {

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
	return &TelegramClient{logger: logger, client: client, cfg: cfg, recvMsgChan: recvMsgChan}, nil
}

func (tc *TelegramClient) Connect() error {

	dispatcher := ext.NewDispatcher(&ext.DispatcherOpts{})
	updates := ext.NewUpdater(&ext.UpdaterOpts{Dispatcher: dispatcher})

	updates.StartPolling(tc.client, &ext.PollingOpts{
		DropPendingUpdates: true,
	})

	handler := handler.NewHandler(tc.logger, tc.recvMsgChan)
	dispatcher.AddHandler(handler)
	return nil
}

func (tc *TelegramClient) Disconnect() {
}

func (tc *TelegramClient) SendMessage(msg *ma.Message) error {
	if msg == nil {
		return nil
	}

	switch msg.Type {
	case ma.TEXT:
		opts := &tgbotapi.SendMessageOpts{ReplyMarkup: msg.GetTgMarkup()}
		_, err := tc.client.SendMessage(msg.GetTgID(), msg.Text, opts)
		return err
	case ma.IMAGE:
		opts := &tgbotapi.SendPhotoOpts{Caption: msg.Text}
		_, err := tc.client.SendPhoto(msg.GetTgID(), msg.Image, opts)
		return err
	}
	return nil
}

func (tc *TelegramClient) GetType() ma.MessageSource {
	return ma.TELEGRAM
}

func (tc *TelegramClient) DownloadFile(id string, msg *ma.Message) string {
	length := len(msg.Data.TgData.Photo)
	if length == 0 {
		return ""
	}

	photo := msg.Data.TgData.Photo[length-1]
	file, err := tc.client.GetFile(photo.FileId, nil)
	if err != nil {
		return ""
	}

	tc.logger.Infof("Dwonloading file %s", file.GetURL(tc.client))
	resp, err := http.Get(file.GetURL(tc.client))
	if err != nil {
		tc.logger.Error(err)
		return ""
	}
	defer resp.Body.Close()

	if err := os.MkdirAll(fmt.Sprintf("./webapp/masters/images/%s/", id), os.ModePerm); err != nil {
		tc.logger.Error(err)
		return ""
	}

	path := fmt.Sprintf("./webapp/masters/images/%s/%s.jpeg", id, file.FileId)
	out, err := os.Create(path)
	if err != nil {
		tc.logger.Error(err)
		return ""
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		tc.logger.Error(err)
	}
	return path
}
