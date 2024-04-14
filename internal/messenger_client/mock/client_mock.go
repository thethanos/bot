package messenger_client

import (
	ma "bot/internal/msgadapter"
	"fmt"
	"sync"
)

type ClientMock struct {
	Result chan *ma.Message
	Wg     *sync.WaitGroup
}

func (c *ClientMock) Connect() error {
	return nil
}

func (c *ClientMock) Disconnect() {

}

func (c *ClientMock) SendMessage(msg *ma.Message) error {
	fmt.Println(msg.Text)
	c.Wg.Done()
	return nil
}

func (c *ClientMock) GetType() ma.MessageSource {
	return ma.TELEGRAM
}

func (c *ClientMock) DownloadFile(ma.FileType, *ma.Message) []byte {
	return []byte{}
}
