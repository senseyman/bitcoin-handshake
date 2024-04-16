package service

import (
	"net"
	"time"

	"github.com/senseyman/bitcoin-handshake/model"
)

type MessageGenerator struct {
}

func NewMessageGenerator() *MessageGenerator {
	return &MessageGenerator{}
}

func (g *MessageGenerator) GenerateNewVersionMessage(
	remoteHost string, remotePort uint16,
	localHost string, localPort uint16,
) model.VersionMessage {
	return model.VersionMessage{
		Version:   model.ProtocolVersion,
		Services:  model.ServiceNodeNetwork,
		Timestamp: time.Now().Unix(),
		AddrRecv: model.NetAddress{
			Timestamp: time.Now().Unix(),
			Services:  model.ServiceNodeNetwork,
			IP:        net.ParseIP(remoteHost),
			Port:      remotePort,
		},
		AddrFrom: model.NetAddress{
			Timestamp: time.Now().Unix(),
			Services:  model.ServiceNodeNetwork,
			IP:        net.ParseIP(localHost),
			Port:      localPort,
		},
		Nonce:       uint64(time.Now().Unix()),
		UserAgent:   "/sensei:0.0.1/",
		StartHeight: 0,
		Relay:       true,
	}
}
