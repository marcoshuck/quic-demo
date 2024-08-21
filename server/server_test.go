package server

import (
	"github.com/marcoshuck/quic-demo/camera"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net"
	"testing"
)

func TestNewServer(t *testing.T) {
	var conn net.PacketConn
	srv := NewQuicServer(conn, camera.NewCamera(nil), nil, nil)
	assert.Implements(t, (*Server)(nil), srv)
}

type testServer struct {
	mock.Mock
}

func (t *testServer) Listen() error {
	args := t.Called()
	return args.Error(0)
}
