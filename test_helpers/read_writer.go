package test_helpers

import (
	"bytes"
	"net"
	"time"
)

var zeroTCPAddr = &net.TCPAddr{
	IP: net.IPv4zero,
}

type ConnWrap struct {
	net.Conn
	R bytes.Buffer
	W bytes.Buffer
}

// Ovrds for some net.Conn methods

func (cw *ConnWrap) Close() error {
	return nil
}

func (cw *ConnWrap) Read(b []byte) (int, error) {
	return cw.R.Read(b)
}

func (cw *ConnWrap) Write(b []byte) (int, error) {
	return cw.W.Write(b)
}

func (cw *ConnWrap) RemoteAddr() net.Addr {
	return zeroTCPAddr
}

func (cw *ConnWrap) LocalAddr() net.Addr {
	return zeroTCPAddr
}

func (cw *ConnWrap) SetDeadline(_ time.Time) error {
	return nil
}

func (cw *ConnWrap) SetReadDeadline(_ time.Time) error {
	return nil
}

func (cw *ConnWrap) SetWriteDeadline(_ time.Time) error {
	return nil
}
