package core

import (
	"sync"

	"github.com/senseyman/bitcoin-handshake/model"
)

const (
	receiveChannelSize = 10
)

type Core struct {
	messageReceiveOnce sync.Once
	connectOnce        sync.Once
	decoder            Decoder
	encoder            Encoder
	generator          Generator
	client             Client

	receiveCh chan model.MessageFromNode
}

func New(decoder Decoder, encoder Encoder, generator Generator, client Client) *Core {
	c := &Core{
		messageReceiveOnce: sync.Once{},
		connectOnce:        sync.Once{},
		decoder:            decoder,
		encoder:            encoder,
		generator:          generator,
		client:             client,
		receiveCh:          make(chan model.MessageFromNode, receiveChannelSize),
	}

	return c
}

func (c *Core) GetReceiveChannel() chan model.MessageFromNode {
	return c.receiveCh
}
