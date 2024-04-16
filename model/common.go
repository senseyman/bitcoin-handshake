package model

import (
	"net"
)

type NetAddress struct {
	Timestamp int64
	Services  uint64
	IP        net.IP
	Port      uint16
}

type MessageHeader struct {
	Magic    uint32  // 4 bytes
	Command  string  // 12 bytes
	Length   uint32  // 4 bytes
	Checksum [4]byte // 4 bytes
}

type EmptyMessage struct {
}

type MessageFromNode struct {
	Header  MessageHeader
	Payload any

	Error *error
}
