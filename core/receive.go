package core

import (
	"bytes"
	"context"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/senseyman/bitcoin-handshake/model"
)

func (c *Core) ReceiveMessages(ctx context.Context) {
	c.messageReceiveOnce.Do(func() {
		go c.client.ReceiveMsg(ctx, c.readHeader, c.payloadRead, c.receiveCh)
	})
}

func (c *Core) readHeader(reader *bytes.Reader) (model.MessageHeader, error) {
	var (
		hdr     model.MessageHeader
		command [model.CommandSize]byte
	)
	err := c.decoder.DecodeElements(reader, &hdr.Magic, &command, &hdr.Length, &hdr.Checksum)
	hdr.Command = string(bytes.TrimRight(command[:], "\x00"))

	return hdr, err
}
func (c *Core) payloadRead(reader *bytes.Reader, header model.MessageHeader) (any, error) {
	switch header.Command {
	case model.VersionCMD:
		versionMsg := model.VersionMessage{}
		err := c.decoder.DecodeElements(
			reader,
			&versionMsg.Version,
			&versionMsg.Services,
			&versionMsg.Timestamp,
			&versionMsg.AddrFrom,
			&versionMsg.AddrRecv,
			&versionMsg.Nonce,
			&versionMsg.UserAgent,
			&versionMsg.StartHeight,
			&versionMsg.Relay,
		)
		return versionMsg, err
	case model.VerackCMD:
		return model.EmptyMessage{}, nil
	}

	return nil, fmt.Errorf("unknown command, can't parse payload: %s", header.Command)
}

func (c *Core) listenReceiveChannel(ctx context.Context, versionMsgLockCh, verackMsgLockCh chan struct{}) {
	receiveCh := c.GetReceiveChannel()

	log.Info("Starting listening incoming messages from node...")
	for {
		select {
		case msg := <-receiveCh:
			if msg.Error != nil {
				log.Errorf("got invalid message: %v", *msg.Error)
			}
			if msg.Header.Command == model.VersionCMD {
				log.Info("got version message")
				log.Infof("%+v\n", msg)
				versionMsgLockCh <- struct{}{}
			}
			if msg.Header.Command == model.VerackCMD {
				log.Info("got verack message")
				log.Infof("%+v\n", msg)
				verackMsgLockCh <- struct{}{}
				return
			}
			log.Info(msg)
		case <-ctx.Done():
			log.Warn("stopping listening messages from node by timeout")
			return
		}
	}
}

func (c *Core) Handshake(ctx context.Context) (int64, error) {
	versionMsgLockCh := make(chan struct{})
	verackMsgLockCh := make(chan struct{})

	// go routing for processing messages from node
	go c.listenReceiveChannel(ctx, versionMsgLockCh, verackMsgLockCh)

	handshakeStartTime := time.Now()
	err := c.sendHandshakeMessages(ctx, versionMsgLockCh, verackMsgLockCh)
	if err != nil {
		return 0, err
	}

	return time.Since(handshakeStartTime).Milliseconds(), nil
}

func (c *Core) sendHandshakeMessages(ctx context.Context, versionMsgLockCh, verackMsgLockCh chan struct{}) error {
	// sending version message to node. This it the first mandatory message we need to send to start our handshake process
	if err := c.SendVersionMessage(); err != nil {
		log.Errorf("err sending version message to node: %v", err)
		return err
	}
	select {
	// if we receive version message from node after our one, we can continue with sending verack message
	case <-versionMsgLockCh:
		log.Info("version message received successfully, trying to send verack message")
	case <-ctx.Done():
		log.Warn("stopping sending version message by context cancel")
		return model.ErrContextTimeout
	}

	// sending verack message to node. This it the second mandatory message we need to send to start our handshake process
	if err := c.SendVerackMessage(); err != nil {
		log.Errorf("err sending verack message to node: %v", err)
		return err
	}

	select {
	// if we receive verack message from node after our one, our handshake is finished
	case <-verackMsgLockCh:
		log.Info("verack message received successfully")
	case <-ctx.Done():
		log.Warn("stopping sending verack message by context cancel")
		return model.ErrContextTimeout
	}

	return nil
}
