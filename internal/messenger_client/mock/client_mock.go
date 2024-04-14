package messenger_client

import ma "bot/internal/msgadapter"

type ClientMock struct {
}

func (c *ClientMock) Connect() error {
	return nil
}

func (c *ClientMock) Disconnect() {

}

func (c *ClientMock) SendMessage(*ma.Message) error {
	return nil
}

func (c *ClientMock) GetType() ma.MessageSource {
	return ma.TELEGRAM
}

func (c *ClientMock) DownloadFile(ma.FileType, *ma.Message) []byte {
	return []byte{}
}
