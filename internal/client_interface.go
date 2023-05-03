package client_interface

type ClientInterface interface {
	Connect() error
	Disconnect()
	SendMessage(message string)
}
