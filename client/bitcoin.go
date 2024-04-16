package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/senseyman/bitcoin-handshake/model"
	"github.com/senseyman/bitcoin-handshake/utils"
)

type BitcoinClient struct {
	conn Connection

	nodeHost string
	nodePort int

	connectionFn func(host string, port int) (Connection, error)

	isConnected bool
}

func NewBitcoinClient(host string, port int,
	connectionFn func(host string, port int) (Connection, error)) (*BitcoinClient, error) {
	b := &BitcoinClient{
		nodeHost:     host,
		nodePort:     port,
		connectionFn: connectionFn,
	}

	if err := b.connect(); err != nil {
		return nil, err
	}

	return b, nil
}

func (c *BitcoinClient) connect() error {
	log.Debug("connecting to node...")
	c.isConnected = false
	conn, err := c.connectionFn(c.nodeHost, c.nodePort)
	if err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}

	c.conn = conn
	c.isConnected = true

	return nil
}

func (c *BitcoinClient) GetNodeHost() string {
	return c.nodeHost
}

func (c *BitcoinClient) GetNodePort() int {
	return c.nodePort
}

func (c *BitcoinClient) Write(msg []byte) (n int, err error) {
	if !c.isConnected {
		return 0, model.ErrConnectionClosed
	}
	return c.conn.Write(msg)
}

func (c *BitcoinClient) ReceiveMsg(
	ctx context.Context,
	headerReadFn func(reader *bytes.Reader) (model.MessageHeader, error),
	payloadReadFn func(reader *bytes.Reader, header model.MessageHeader) (any, error),
	receiveCh chan model.MessageFromNode,
) {
	for {
		select {
		case <-ctx.Done():
			log.Warn("stopping receiving thread by context done")
			if err := c.conn.Close(); err != nil {
				log.Warnf("err while closing connection to node: %v", err)
			}
			c.isConnected = false
			return
		case <-time.NewTicker(time.Millisecond).C:
			c.receive(headerReadFn, payloadReadFn, receiveCh)
		}
	}
}

func (c *BitcoinClient) receive(
	headerReadFn func(reader *bytes.Reader) (model.MessageHeader, error),
	payloadReadFn func(reader *bytes.Reader, header model.MessageHeader) (any, error),
	receiveCh chan model.MessageFromNode,
) {
	if !c.isConnected {
		if err := c.connect(); err != nil {
			log.Errorf("err while reconnectiong to blockchain node: %v", err)
		}
		return
	}

	// read header to determine message type and payload size
	var headerBytes [24]byte
	_, err := io.ReadFull(c.conn, headerBytes[:])
	if err != nil {
		if errors.Is(err, io.EOF) {
			c.isConnected = false
			log.Debug("got EOF")
		}
		return
	}
	hr := bytes.NewReader(headerBytes[:])
	hdr, err := headerReadFn(hr)
	if err != nil {
		log.Warnf("err while parsing msg header: %v", err)
		return
	}

	// validate magic number
	if hdr.Magic != model.TestNetMagic {
		log.Warnf("got mesage with invalid magic number")
		receiveCh <- model.MessageFromNode{
			Error: &model.ErrInvalidMagicNumber,
		}
		return
	}

	payloadBytes := make([]byte, hdr.Length)
	_, err = io.ReadFull(c.conn, payloadBytes[:])
	plr := bytes.NewReader(payloadBytes[:])
	if err != nil {
		log.Warnf("err while reading msg payload from the connection: %v", err)
		return
	}

	// check if checksum is valid
	actualChecksum := utils.DoubleHashB(payloadBytes)[0:4]
	if !bytes.Equal(hdr.Checksum[:], actualChecksum) {
		log.Warnf("got mesage with invalid checksum")
		receiveCh <- model.MessageFromNode{
			Error: &model.ErrInvalidMessageChecksum,
		}
		return
	}

	msg, err := payloadReadFn(plr, hdr)
	if err != nil {
		log.Warnf("err while parsing msg payload: %v. Skipping", err)
		return
	}

	log.Debug("sending read message from node to processing")
	receiveCh <- model.MessageFromNode{
		Header:  hdr,
		Payload: msg,
	}
}
