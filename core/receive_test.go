package core

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/senseyman/bitcoin-handshake/core/mock"
	"github.com/senseyman/bitcoin-handshake/model"
)

type failType int

const (
	failNone failType = iota
	failTimeoutVerack
	failTimeoutVersion
	failSendVerack
	failSendVersion
)

var (
	testErr = errors.New("test error")
)

func TestCore_ReceiveMessages(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	client.EXPECT().ReceiveMsg(ctx, gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()

	New(nil, nil, nil, client).ReceiveMessages(ctx)
}

func TestCore_Handshake(t *testing.T) {
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
	)

	testCases := []struct {
		name string
		init func(t *testing.T, fail failType,
			remoteHost string, remotePort int,
			localHost string, locaPort int,
			versionMsg model.VersionMessage) *Core
		hasErr bool
		fail   failType
		expErr error
	}{
		{
			name:   "success",
			fail:   failNone,
			init:   prepareMocks,
			hasErr: false,
			expErr: nil,
		},
		{
			name:   "err/ctx_timeout_verack",
			fail:   failTimeoutVerack,
			init:   prepareMocks,
			hasErr: true,
			expErr: model.ErrContextTimeout,
		},
		{
			name:   "err/ctx_timeout_version",
			fail:   failTimeoutVersion,
			init:   prepareMocks,
			hasErr: true,
			expErr: model.ErrContextTimeout,
		},
		{
			name:   "err/send_verack",
			fail:   failSendVerack,
			init:   prepareMocks,
			hasErr: true,
			expErr: testErr,
		},
		{
			name:   "err/send_version",
			fail:   failSendVersion,
			init:   prepareMocks,
			hasErr: true,
			expErr: testErr,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*6)
			defer cancel()

			c := tc.init(t, tc.fail, remoteHost, remotePort, localHost, locaPort, versionMsg)
			_, err := c.Handshake(ctx)

			if tc.hasErr {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func prepareMocks(t *testing.T, fail failType,
	remoteHost string, remotePort int,
	localHost string, locaPort int,
	versionMsg model.VersionMessage,
) *Core {
	// init mocks
	ctrl := gomock.NewController(t)
	client := mock.NewMockClient(ctrl)
	encoder := mock.NewMockEncoder(ctrl)
	generator := mock.NewMockGenerator(ctrl)

	c := New(nil, encoder, generator, client)
	recCh := c.GetReceiveChannel()

	switch fail {
	case failNone:
		mockSendVersionMessage(client, encoder, generator, remoteHost, remotePort, localHost, locaPort, versionMsg, fail)
		mockSendVerackMessage(client, encoder, fail)
		go func() {
			time.Sleep(time.Second * 2)
			recCh <- model.MessageFromNode{
				Header: model.MessageHeader{Command: model.VersionCMD},
			}
		}()
		go func() {
			time.Sleep(time.Second * 4)
			recCh <- model.MessageFromNode{
				Header: model.MessageHeader{Command: model.VerackCMD},
			}
		}()
	case failTimeoutVerack, failSendVerack:
		mockSendVersionMessage(client, encoder, generator, remoteHost, remotePort, localHost, locaPort, versionMsg, fail)
		mockSendVerackMessage(client, encoder, fail)
		go func() {
			time.Sleep(time.Second * 2)
			recCh <- model.MessageFromNode{
				Header: model.MessageHeader{Command: model.VersionCMD},
			}
		}()
	case failTimeoutVersion, failSendVersion:
		mockSendVersionMessage(client, encoder, generator, remoteHost, remotePort, localHost, locaPort, versionMsg, fail)
	}

	return c
}

func mockSendVersionMessage(
	client *mock.MockClient,
	encoder *mock.MockEncoder,
	generator *mock.MockGenerator,
	remoteHost string, remotePort int,
	localHost string, locaPort int,
	versionMsg model.VersionMessage,
	fail failType,
) {
	client.EXPECT().GetNodeHost().Return(remoteHost)
	client.EXPECT().GetNodePort().Return(remotePort)
	generator.EXPECT().GenerateNewVersionMessage(remoteHost, uint16(remotePort),
		localHost, uint16(locaPort)).Return(versionMsg)
	encoder.EXPECT().EncodeVersionMessage(gomock.Any(), versionMsg).Return(nil)
	encoder.EXPECT().EncodeElements(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	client.EXPECT().Write(gomock.Any()).Return(0, nil)

	if fail == failSendVersion {
		client.EXPECT().Write(gomock.Any()).Return(0, testErr)
	} else {
		client.EXPECT().Write(gomock.Any()).Return(0, nil)
	}
}

func mockSendVerackMessage(
	client *mock.MockClient,
	encoder *mock.MockEncoder,
	fail failType,
) {
	encoder.EXPECT().EncodeElements(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
	if fail == failSendVerack {
		client.EXPECT().Write(gomock.Any()).Return(0, testErr)
	} else {
		client.EXPECT().Write(gomock.Any()).Return(0, nil)
	}
}
