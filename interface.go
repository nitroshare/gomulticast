package gomulticast

import (
	"errors"
	"net"
)

var (
	ErrNoPackets = errors.New("no packets to be received")
)

// Interface wraps the information in net.Interface. This allows for much
// easier mocking.
type Interface interface {

	// Interface returns the underlying net.Interface.
	Interface() *net.Interface
}

type netInterface struct {
	i *net.Interface
}

func (n netInterface) Interface() *net.Interface { return n.i }

// Packet represents a packet sent or received on a MockInterface.
type Packet struct {
	Addr net.Addr
	Data []byte
}

// MockInterface implements Interface and provides methods for sending and
// receiving packets on the interface.
type MockInterface struct {
	sendQueue []*Packet
	recvQueue []*Packet
}

func (m MockInterface) Interface() *net.Interface { return nil }

// Send queues the provided packet from the provided address for being
// received by a listener.
func (m *MockInterface) Send(p *Packet) {
	m.sendQueue = append(m.sendQueue, p)
}

// Receive dequeues a packet that was sent by a listener. An error is returned
// if there were no packets sent.
func (m *MockInterface) Receive() (*Packet, error) {
	if len(m.recvQueue) == 0 {
		return nil, ErrNoPackets
	}
	p := m.recvQueue[0]
	m.recvQueue = m.recvQueue[1:]
	return p, nil
}

var (
	// Interfaces returns a list of interfaces on the host when not mocked.
	Interfaces func() ([]Interface, error)

	origNetInterfaces = net.Interfaces
)

func netInterfaces() ([]Interface, error) {
	l, err := origNetInterfaces()
	if err != nil {
		return nil, err
	}
	interfaces := []Interface{}
	for _, i := range l {
		interfaces = append(interfaces, &netInterface{&i})
	}
	return interfaces, nil
}

var (
	mInterfaces []Interface
)

func mockInterfaces() ([]Interface, error) {
	return mInterfaces, nil
}

// AddMockInterface adds a MockInterface to the list returned by Interfaces
// when mocking is enabled.
func AddMockInterface(i *MockInterface) {
	mInterfaces = append(mInterfaces, i)
}
