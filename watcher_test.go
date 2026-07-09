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
	var (
		chanAdded   = make(chan Interface)
		chanRemoved = make(chan Interface)
		w           = NewWatcher(&WatcherConfig{
			Interval:    time.Second,
			ChanAdded:   chanAdded,
			ChanRemoved: chanRemoved,
		})
	)
	defer w.Close()
	AddMockInterface(NewMockInterface())
	<-chanAdded
	mInterfaces = []Interface{}
	gotime.Advance(2 * time.Second)
	<-chanRemoved
	gotime.Advance(2 * time.Second)
	select {
	case <-chanAdded:
		t.Fatal("unexpected interface added")
	case <-chanRemoved:
		t.Fatal("unexpected interface removed")
	default:
	}
}

func TestWatcherError(t *testing.T) {
	origNetInterfaces = func() ([]net.Interface, error) { return nil, testError }
	defer func() {
		origNetInterfaces = net.Interfaces
	}()
	var (
		chanAdded   = make(chan Interface)
		chanRemoved = make(chan Interface)
		w           = NewWatcher(&WatcherConfig{
			Interval:    time.Second,
			ChanAdded:   chanAdded,
			ChanRemoved: chanRemoved,
		})
	)
	defer w.Close()
}
