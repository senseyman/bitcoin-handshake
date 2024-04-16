package core

import (
	"bytes"

	log "github.com/sirupsen/logrus"

	"github.com/senseyman/bitcoin-handshake/model"
	"github.com/senseyman/bitcoin-handshake/utils"
)

func (c *Core) SendVersionMessage() error {
	log.Info("sending version message")
	var command [model.CommandSize]byte
	copy(command[:], model.VersionCMD)

	msg := c.generator.GenerateNewVersionMessage(
		c.client.GetNodeHost(), uint16(c.client.GetNodePort()),
		"127.0.0.1", 0, // ignore it as node will respond us with our correct white IP
	)
	var bw bytes.Buffer
	if err := c.encoder.EncodeVersionMessage(&bw, msg); err != nil {
		return err
	}

	payload := bw.Bytes()
	payloadLen := len(payload)

	hdr := model.MessageHeader{}
	hdr.Magic = model.TestNetMagic
	hdr.Command = model.VersionCMD
	hdr.Length = uint32(payloadLen)
	copy(hdr.Checksum[:], utils.DoubleHashB(payload)[0:4])

	hw := bytes.NewBuffer(make([]byte, 0, 24))
	if err := c.encoder.EncodeElements(hw, hdr.Magic, command, hdr.Length, hdr.Checksum); err != nil {
		return err
	}

	n, err := c.client.Write(hw.Bytes())
	if err != nil {
		return err
	}
	log.Debugf("sent header, number of byte %d", n)

	n, err = c.client.Write(payload)
	if err != nil {
		return err
	}
	log.Debugf("sent payload, number of byte %d", n)

	return nil
}

func (c *Core) SendVerackMessage() error {
	log.Info("sending verack message")
	var command [model.CommandSize]byte
	copy(command[:], model.VerackCMD)

	hdr := model.MessageHeader{}
	hdr.Magic = model.TestNetMagic
	hdr.Command = model.VerackCMD

	hw := bytes.NewBuffer(make([]byte, 0, 24))
	if err := c.encoder.EncodeElements(hw, hdr.Magic, command, hdr.Length, hdr.Checksum); err != nil {
		return err
	}

	n, err := c.client.Write(hw.Bytes())
	if err != nil {
		return err
	}
	log.Debugf("sent header, number of byte %d", n)

	return nil
}
