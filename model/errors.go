package model

import (
	"errors"
)

var (
	ErrContextTimeout         = errors.New("cancel by context timeout")
	ErrConnectionClosed       = errors.New("connection to node is closed")
	ErrInvalidMessageChecksum = errors.New("invalid message checksum")
	ErrInvalidMagicNumber     = errors.New("invalid message magic number")
)
