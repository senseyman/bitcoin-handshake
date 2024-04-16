package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	assert.NotNil(t, New(nil, nil, nil, nil))
}

func TestCore_GetReceiveChannel(t *testing.T) {
	c := New(nil, nil, nil, nil)
	recCh := c.GetReceiveChannel()

	assert.NotNil(t, recCh)
	assert.Equal(t, receiveChannelSize, cap(recCh))
}
