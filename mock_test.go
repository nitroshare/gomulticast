package gomulticast

import (
	"testing"

	"github.com/nitroshare/compare"
)

func TestMockUnmock(t *testing.T) {
	compare.CompareFn(t, Interfaces, netInterfaces, true)
	compare.CompareFn(t, listenMulticastUDP, netListenMulticastUDP, true)
	Mock()
	compare.CompareFn(t, Interfaces, mockInterfaces, true)
	compare.CompareFn(t, listenMulticastUDP, mockListenMulticastUDP, true)
	Unmock()
	compare.CompareFn(t, Interfaces, netInterfaces, true)
	compare.CompareFn(t, listenMulticastUDP, netListenMulticastUDP, true)
}
