package gomulticast

import (
	"net"
	"testing"

	"github.com/nitroshare/compare"
)

var (
	testAddr net.Addr = nil
	testData          = []byte("test")
)

func TestMockInterfaceInterface(t *testing.T) {
	NewMockInterface().Interface()
}

func TestMockInterfaceQueueForRead(t *testing.T) {
	var (
		m         = NewMockInterface()
		c         = &mockUDPConn{m}
		chanClose = make(chan any)
		errChan   = make(chan error)
		b         = make([]byte, 32)
	)
	m.chanTest = make(chan any)
	go func() {
		defer close(chanClose)
		c.ReadFrom(b)
	}()
	<-m.chanTest
	m.QueueForRead(&Packet{testAddr, testData})
	<-chanClose
	m.QueueForRead(&Packet{testAddr, testData})
	n, addr, err := c.ReadFrom(b)
	compare.Compare(t, string(b[:n]), string(testData), true)
	compare.Compare(t, n, len(testData), true)
	compare.Compare(t, addr, testAddr, true)
	compare.Compare(t, err, nil, true)
	go func() {
		_, _, err := c.ReadFrom(b)
		errChan <- err
	}()
	<-m.chanTest
	c.Close()
	v := <-errChan
	compare.Compare(t, v, net.ErrClosed, true)
}

func TestMockInterfaceDequeueWrite(t *testing.T) {
	var (
		m = NewMockInterface()
		c = &mockUDPConn{m}
	)
	c.WriteTo(testData, testAddr)
	_, err := m.DequeueWrite()
	compare.Compare(t, err, nil, true)
	var (
		chanClose = make(chan any)
		p         *Packet
		pErr      error
	)
	m.chanTest = make(chan any)
	go func() {
		defer close(chanClose)
		p, pErr = m.DequeueWrite()
	}()
	<-m.chanTest
	n, err := c.WriteTo(testData, testAddr)
	compare.Compare(t, n, len(testData), true)
	compare.Compare(t, err, nil, true)
	<-chanClose
	compare.Compare(t, p.Addr, testAddr, true)
	compare.Compare(t, string(p.Data), string(testData), true)
	compare.Compare(t, pErr, nil, true)
	go func() {
		p, err = m.DequeueWrite()
	}()
	<-m.chanTest
	c.Close()
	compare.Compare(t, p, nil, true)
	compare.Compare(t, err, net.ErrClosed, true)
}
