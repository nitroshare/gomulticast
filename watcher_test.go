package gomulticast

import (
	"net"
	"testing"
	"time"

	"github.com/nitroshare/gotime"
)

func TestWatcher(t *testing.T) {
	gotime.Mock()
	defer gotime.Unmock()
	Mock()
	defer Unmock()
	AddMockInterface(NewMockInterface())
	w := NewWatcher(time.Second)
	<-w.ChanAdded
	mInterfaces = []Interface{}
	gotime.Advance(2 * time.Second)
	<-w.ChanRemoved
	gotime.Advance(2 * time.Second)
	select {
	case <-w.ChanAdded:
		t.Fatal("unexpected interface added")
	case <-w.ChanRemoved:
		t.Fatal("unexpected interface removed")
	default:
	}
	defer w.Close()
}

func TestWatcherError(t *testing.T) {
	origNetInterfaces = func() ([]net.Interface, error) { return nil, testError }
	defer func() {
		origNetInterfaces = net.Interfaces
	}()
	w := NewWatcher(time.Second)
	defer w.Close()
}
