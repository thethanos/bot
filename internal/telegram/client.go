package telegram

import (
	"bot/internal/config"
	"bot/internal/logger"
	ma "bot/internal/msgadapter"
	"io"
	"net/http"

	handler "bot/internal/telegram/event_handler"

	"github.com/PaulSonOfLars/gotgbot/v2"
	tgbotapi "github.com/PaulSonOfLars/gotgbot/v2"
	"github.com/PaulSonOfLars/gotgbot/v2/ext"
)

type TelegramClient struct {
	logger      logger.Logger
	cfg         *config.Config
	recvMsgChan chan *ma.Message
	client      *tgbotapi.Bot
}

func NewTelegramClient(logger logger.Logger, cfg *config.Config, recvMsgChan chan *ma.Message) (*TelegramClient, error) {

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

	if err := updates.StartPolling(tc.client, &ext.PollingOpts{
		DropPendingUpdates: true,
	}); err != nil {
		return err
	}

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

func (tc *TelegramClient) DownloadFile(fileType ma.FileType, msg *ma.Message) []byte {

	var fileId string
	switch fileType {
	case ma.DOCUMENT:
		if msg.Data.TgData.Document == nil {
			tc.logger.Error("empty document")
			return nil
		}
		fileId = msg.Data.TgData.Document.FileId
	case ma.PHOTO:
		if msg.Data.TgData.Photo == nil {
			tc.logger.Error("empty picture")
			return nil
		}
		length := len(msg.Data.TgData.Photo)
		fileId = msg.Data.TgData.Photo[length-1].FileId
	default:
		tc.logger.Error("file type is not supported")
		return nil
	}

	file, err := tc.client.GetFile(fileId, nil)
	if err != nil {
		return nil
	}

	tc.logger.Infof("Dwonloading file %s", file.GetURL(tc.client))
	resp, err := http.Get(file.GetURL(tc.client))
	if err != nil {
		tc.logger.Error(err)
		return nil
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	return data
}
