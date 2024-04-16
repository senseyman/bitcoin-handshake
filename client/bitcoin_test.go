package client

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/senseyman/bitcoin-handshake/client/mock"
)

func TestNewBitcoinClient(t *testing.T) {
	testCases := []struct {
		name   string
		init   func(t *testing.T) (*BitcoinClient, error)
		hasErr bool
	}{
		{
			name: "success",
			init: func(t *testing.T) (*BitcoinClient, error) {
				ctrl := gomock.NewController(t)
				conn := mock.NewMockConnection(ctrl)

				fn := func(host string, port int) (Connection, error) {
					return conn, nil
				}

				return NewBitcoinClient("127.0.0.1", 8333, fn)
			},
			hasErr: false,
		},
		{
			name: "err/connection",
			init: func(t *testing.T) (*BitcoinClient, error) {
				fn := func(host string, port int) (Connection, error) {
					return nil, errors.New("test error")
				}

				return NewBitcoinClient("127.0.0.1", 8333, fn)
			},
			hasErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cli, err := tc.init(t)
			if tc.hasErr {
				assert.Error(t, err)
				assert.Nil(t, cli)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, cli)
			}
		})
	}
}

func TestBitcoinClient_Write(t *testing.T) {
	host := "127.0.0.1"
	port := 8333
	testErr := errors.New("test error")

	testCases := []struct {
		name   string
		init   func(t *testing.T, msg []byte) (*BitcoinClient, error)
		hasErr bool
	}{
		{
			name: "success",
			init: func(t *testing.T, msg []byte) (*BitcoinClient, error) {
				ctrl := gomock.NewController(t)
				conn := mock.NewMockConnection(ctrl)

				fn := func(host string, port int) (Connection, error) {
					return conn, nil
				}

				conn.EXPECT().Write(msg).Return(len(msg), nil)

				return NewBitcoinClient(host, port, fn)
			},
			hasErr: false,
		},
		{
			name: "err/write",
			init: func(t *testing.T, msg []byte) (*BitcoinClient, error) {
				ctrl := gomock.NewController(t)
				conn := mock.NewMockConnection(ctrl)

				fn := func(host string, port int) (Connection, error) {
					return conn, nil
				}

				conn.EXPECT().Write(msg).Return(0, testErr)

				return NewBitcoinClient(host, port, fn)
			},
			hasErr: true,
		},
		{
			name: "err/not_connected",
			init: func(t *testing.T, msg []byte) (*BitcoinClient, error) {
				ctrl := gomock.NewController(t)
				conn := mock.NewMockConnection(ctrl)

				fn := func(host string, port int) (Connection, error) {
					return conn, nil
				}

				c, _ := NewBitcoinClient(host, port, fn)
				c.isConnected = false

				return c, nil
			},
			hasErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			msg := []byte{1, 2, 3, 4, 5, 6}
			c, _ := tc.init(t, msg)

			n, err := c.Write(msg)
			if tc.hasErr {
				assert.Zero(t, n)
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(msg), n)
			}
		})
	}
}
