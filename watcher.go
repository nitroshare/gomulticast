package gomulticast

import (
	"time"

	"github.com/nitroshare/gotime"
)

// Watcher monitors available network interfaces and notifies when one is
// added or removed.
type Watcher struct {
	// ChanAdded receives an Interface when a new interface is added.
	ChanAdded <-chan Interface
	// ChanRemoved receives an Interface when an interface is removed.
	ChanRemoved <-chan Interface

	chanAdded   chan<- Interface
	chanRemoved chan<- Interface
	chanClose   chan any
	chanClosed  chan any
}

func (w *Watcher) diff(m map[string]Interface) map[string]Interface {
	interfaces, err := Interfaces()
	if err != nil {
		return m
	}
	m2 := map[string]Interface{}
	for _, i := range interfaces {
		m2[i.Interface().Name] = i
	}
	for k, v := range m {
		if _, ok := m2[k]; !ok {
			w.chanRemoved <- v
		}
	}
	for k, v := range m2 {
		if _, ok := m[k]; !ok {
			w.chanAdded <- v
		}
	}
	return m2
}

func (w *Watcher) run(interval time.Duration) {
	defer close(w.chanClosed)
	defer close(w.chanRemoved)
	defer close(w.chanAdded)
	var (
		m = w.diff(map[string]Interface{})
		t = gotime.NewTicker(interval)
	)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			m = w.diff(m)
		case <-w.chanClose:
			return
		}
	}
}

// NewWatcher creates a new Watcher instance using the provided interval.
func NewWatcher(interval time.Duration) *Watcher {
	var (
		chanAdded   = make(chan Interface)
		chanRemoved = make(chan Interface)
		w           = &Watcher{
			ChanAdded:   chanAdded,
			ChanRemoved: chanRemoved,
			chanAdded:   chanAdded,
			chanRemoved: chanRemoved,
			chanClose:   make(chan any),
			chanClosed:  make(chan any),
		}
	)
	go w.run(interval)
	return w
}

// Close shuts down the watcher.
func (w *Watcher) Close() {
	close(w.chanClose)
	<-w.chanClosed
}
