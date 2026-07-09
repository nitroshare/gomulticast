package gomulticast

import (
	"net"
	"testing"
	"time"

	"github.com/nitroshare/gotime"
)

type watcherSet struct {
	chanAdded   <-chan Interface
	chanRemoved <-chan Interface
	watcher     *Watcher
}

func newWatcherSet() *watcherSet {
	var (
		chanAdded   = make(chan Interface)
		chanRemoved = make(chan Interface)
	)
	return &watcherSet{
		chanAdded:   chanAdded,
		chanRemoved: chanRemoved,
		watcher: NewWatcher(&WatcherConfig{
			Interval:    time.Second,
			ChanAdded:   chanAdded,
			ChanRemoved: chanRemoved,
		}),
	}
}

func TestWatcher(t *testing.T) {
	gotime.Mock()
	defer gotime.Unmock()
	Mock()
	defer Unmock()
	AddMockInterface(NewMockInterface())
	s := newWatcherSet()
	defer s.watcher.Close()
	<-s.chanAdded
	mInterfaces = []Interface{}
	gotime.Advance(2 * time.Second)
	<-s.chanRemoved
	gotime.Advance(2 * time.Second)
	select {
	case <-s.chanAdded:
		t.Fatal("unexpected interface added")
	case <-s.chanRemoved:
		t.Fatal("unexpected interface removed")
	default:
	}
}

func TestWatcherError(t *testing.T) {
	origNetInterfaces = func() ([]net.Interface, error) { return nil, testError }
	defer func() {
		origNetInterfaces = net.Interfaces
	}()
	s := newWatcherSet()
	defer s.watcher.Close()
}

func TestCloseDuringSend(t *testing.T) {
	gotime.Mock()
	defer gotime.Unmock()
	Mock()
	defer Unmock()
	AddMockInterface(NewMockInterface())

	t.Run("chanAdded", func(t *testing.T) {
		s := newWatcherSet()
		defer s.watcher.Close()
	})

	t.Run("chanRemoved", func(t *testing.T) {
		s := newWatcherSet()
		defer s.watcher.Close()
		<-s.chanAdded
		mInterfaces = []Interface{}
		s.watcher.chanTest = make(chan any)
		gotime.Advance(2 * time.Second)
		<-s.watcher.chanTest
	})
}
