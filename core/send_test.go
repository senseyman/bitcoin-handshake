package core

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/senseyman/bitcoin-handshake/core/mock"
	"github.com/senseyman/bitcoin-handshake/model"
)

func TestCore_SendVersionMessage(t *testing.T) {
	var (
		localHost  = "127.0.0.1"
		locaPort   = 0
		remoteHost = "127.0.0.1"
		remotePort = 8333

		versionMsg = model.VersionMessage{
			Version:   model.ProtocolVersion,
			Services:  0,
			Timestamp: time.Now().Unix(),
			AddrRecv: model.NetAddress{
				IP:   net.ParseIP(remoteHost),
				Port: uint16(remotePort),
			},
			AddrFrom: model.NetAddress{
				IP:   net.ParseIP(localHost),
				Port: uint16(locaPort),
			},
			Nonce:       0,
			UserAgent:   "test-agent/1",
			StartHeight: 0,
			Relay:       false,
		}
		testErr = errors.New("test error")
	)

	testCases := []struct {
		name   string
		init   func(t *testing.T) *Core
		hasErr bool
	}{
		{
			name: "success",
			init: func(t *testing.T) *Core {
				ctrl := gomock.NewController(t)
				client := mock.NewMockClient(ctrl)
				encoder := mock.NewMockEncoder(ctrl)
				generator := mock.NewMockGenerator(ctrl)

				client.EXPECT().GetNodeHost().Return(remoteHost)
				client.EXPECT().GetNodePort().Return(remotePort)
				generator.EXPECT().GenerateNewVersionMessage(remoteHost, uint16(remotePort),
					localHost, uint16(locaPort)).Return(versionMsg)
				encoder.EXPECT().EncodeVersionMessage(gomock.Any(), versionMsg).Return(nil)
				encoder.EXPECT().EncodeElements(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				client.EXPECT().Write(gomock.Any()).Return(0, nil)
				client.EXPECT().Write(gomock.Any()).Return(0, nil)

				return New(nil, encoder, generator, client)
			},
			hasErr: false,
		},
		{
			name: "err/write_payload",
			init: func(t *testing.T) *Core {
				ctrl := gomock.NewController(t)
				client := mock.NewMockClient(ctrl)
				encoder := mock.NewMockEncoder(ctrl)
				generator := mock.NewMockGenerator(ctrl)

				client.EXPECT().GetNodeHost().Return(remoteHost)
				client.EXPECT().GetNodePort().Return(remotePort)
				generator.EXPECT().GenerateNewVersionMessage(remoteHost, uint16(remotePort),
					localHost, uint16(locaPort)).Return(versionMsg)
				encoder.EXPECT().EncodeVersionMessage(gomock.Any(), versionMsg).Return(nil)
				encoder.EXPECT().EncodeElements(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				client.EXPECT().Write(gomock.Any()).Return(0, nil)
				client.EXPECT().Write(gomock.Any()).Return(0, testErr)

				return New(nil, encoder, generator, client)
			},
			hasErr: true,
		},
		{
			name: "err/write_header",
			init: func(t *testing.T) *Core {
				ctrl := gomock.NewController(t)
				client := mock.NewMockClient(ctrl)
				encoder := mock.NewMockEncoder(ctrl)
				generator := mock.NewMockGenerator(ctrl)

				client.EXPECT().GetNodeHost().Return(remoteHost)
				client.EXPECT().GetNodePort().Return(remotePort)
				generator.EXPECT().GenerateNewVersionMessage(remoteHost, uint16(remotePort),
					localHost, uint16(locaPort)).Return(versionMsg)
				encoder.EXPECT().EncodeVersionMessage(gomock.Any(), versionMsg).Return(nil)
				encoder.EXPECT().EncodeElements(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				client.EXPECT().Write(gomock.Any()).Return(0, testErr)

				return New(nil, encoder, generator, client)
			},
			hasErr: true,
		},
		{
			name: "err/encode_header",
			init: func(t *testing.T) *Core {
				ctrl := gomock.NewController(t)
				client := mock.NewMockClient(ctrl)
				encoder := mock.NewMockEncoder(ctrl)
				generator := mock.NewMockGenerator(ctrl)

				client.EXPECT().GetNodeHost().Return(remoteHost)
				client.EXPECT().GetNodePort().Return(remotePort)
				generator.EXPECT().GenerateNewVersionMessage(remoteHost, uint16(remotePort),
					localHost, uint16(locaPort)).Return(versionMsg)
				encoder.EXPECT().EncodeVersionMessage(gomock.Any(), versionMsg).Return(nil)
				encoder.EXPECT().EncodeElements(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)

				return New(nil, encoder, generator, client)
			},
			hasErr: true,
		},
		{
			name: "err/encode_version_msg",
			init: func(t *testing.T) *Core {
				ctrl := gomock.NewController(t)
				client := mock.NewMockClient(ctrl)
				encoder := mock.NewMockEncoder(ctrl)
				generator := mock.NewMockGenerator(ctrl)

				client.EXPECT().GetNodeHost().Return(remoteHost)
				client.EXPECT().GetNodePort().Return(remotePort)
				generator.EXPECT().GenerateNewVersionMessage(remoteHost, uint16(remotePort),
					localHost, uint16(locaPort)).Return(versionMsg)
				encoder.EXPECT().EncodeVersionMessage(gomock.Any(), versionMsg).Return(testErr)

				return New(nil, encoder, generator, client)
			},
			hasErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := tc.init(t)
			err := c.SendVersionMessage()

			if tc.hasErr {
				assert.Error(t, err)
				assert.EqualError(t, err, testErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCore_SendVerackMessage(t *testing.T) {
	var (
		testErr = errors.New("test error")
	)

	testCases := []struct {
		name   string
		init   func(t *testing.T) *Core
		hasErr bool
	}{
		{
			name: "success",
			init: func(t *testing.T) *Core {
				ctrl := gomock.NewController(t)
				client := mock.NewMockClient(ctrl)
				encoder := mock.NewMockEncoder(ctrl)

				encoder.EXPECT().EncodeElements(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				client.EXPECT().Write(gomock.Any()).Return(0, nil)

				return New(nil, encoder, nil, client)
			},
			hasErr: false,
		},
		{
			name: "err/write_header",
			init: func(t *testing.T) *Core {
				ctrl := gomock.NewController(t)
				client := mock.NewMockClient(ctrl)
				encoder := mock.NewMockEncoder(ctrl)

				encoder.EXPECT().EncodeElements(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				client.EXPECT().Write(gomock.Any()).Return(0, testErr)

				return New(nil, encoder, nil, client)
			},
			hasErr: true,
		},
		{
			name: "err/encode_header",
			init: func(t *testing.T) *Core {
				ctrl := gomock.NewController(t)
				client := mock.NewMockClient(ctrl)
				encoder := mock.NewMockEncoder(ctrl)

				encoder.EXPECT().EncodeElements(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(testErr)

				return New(nil, encoder, nil, client)
			},
			hasErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			c := tc.init(t)
			err := c.SendVerackMessage()

			if tc.hasErr {
				assert.Error(t, err)
				assert.EqualError(t, err, testErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
