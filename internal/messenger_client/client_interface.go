package messenger_client

import ma "bot/internal/msgadapter"

type ClientInterface interface {
	Connect() error
	Disconnect()
	SendMessage(*ma.Message) error
	GetType() ma.MessageSource
	DownloadFile(ma.FileType, *ma.Message) []byte
}
