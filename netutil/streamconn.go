package netutil

import (
	"errors"
	"io"
	"net"
	"time"
)

type streamConn struct {
	io.ReadCloser
	io.WriteCloser
}

// NewStreamConn wraps a ReadCloser and a WriteCloser into
// a structure that satisfies net.Conn. This allows for example
// to run a 9P server directly on stdin and stdout.
func NewStreamConn(r io.ReadCloser, w io.WriteCloser) net.Conn {
	return &streamConn{r, w}
}

type streamAddr int

func (streamAddr) Network() string {
	return "stream"
}
func (s streamAddr) String() string {
	return "stream"
}

func (c *streamConn) Close() (err error) {
	err = c.ReadCloser.Close()
	err1 := c.WriteCloser.Close()
	if err == nil {
		err = err1
	}
	return
}

func (c *streamConn) LocalAddr() net.Addr {
	return streamAddr(0)
}

func (c *streamConn) RemoteAddr() net.Addr {
	return streamAddr(0)
}

var errNotSupported = errors.New("deadline not supported")

func (c *streamConn) SetDeadline(t time.Time) error {
	return errNotSupported
}

func (c *streamConn) SetReadDeadline(t time.Time) error {
	return errNotSupported
}

func (c *streamConn) SetWriteDeadline(t time.Time) error {
	return errNotSupported
}
