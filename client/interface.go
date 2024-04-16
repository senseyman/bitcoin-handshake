package client

import (
	"io"
)

type Connection interface {
	io.Reader
	io.Writer
	io.Closer
}
