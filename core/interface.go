package core

import (
	"bytes"
	"context"
	"io"

	"github.com/senseyman/bitcoin-handshake/model"
)

type Decoder interface {
	DecodeElements(r io.Reader, elements ...any) error
}

type Encoder interface {
	EncodeVersionMessage(w io.Writer, msg model.VersionMessage) error
	EncodeElements(w io.Writer, elements ...any) error
}

type Generator interface {
	GenerateNewVersionMessage(remoteHost string, remotePort uint16, localHost string, localPort uint16) model.VersionMessage
}

type Client interface {
	io.Writer
	ReceiveMsg(
		ctx context.Context,
		headerReadFn func(reader *bytes.Reader) (model.MessageHeader, error),
		payloadReadFn func(reader *bytes.Reader, header model.MessageHeader) (any, error),
		receiveCh chan model.MessageFromNode,
	)
	GetNodeHost() string
	GetNodePort() int
}
