package gomulticast

import (
	"time"

	"github.com/nitroshare/gotime"
)

// WatcherConfig provides configuration for Watcher.
type WatcherConfig struct {

	// Interval specifies how often the adapters on the host should be
	// enumerated.
	Interval time.Duration

	// ChanAdded receives an Interface when a new interface is added. This can
	// be left nil if not desired.
	ChanAdded chan<- Interface

	// ChanRemoved receives an Interface when an interface is removed. This
	// can be left nil if not desired.
	ChanRemoved chan<- Interface
}

// Watcher monitors available network interfaces and notifies when one is
// added or removed.
type Watcher struct {
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
	defer func() {
		if w.chanAdded != nil {
			close(w.chanAdded)
		}
		if w.chanRemoved != nil {
			close(w.chanRemoved)
		}
	}()
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

// NewWatcher creates a new Watcher instance.
func NewWatcher(cfg *WatcherConfig) *Watcher {
	w := &Watcher{
		chanAdded:   cfg.ChanAdded,
		chanRemoved: cfg.ChanRemoved,
		chanClose:   make(chan any),
		chanClosed:  make(chan any),
	}
	go w.run(cfg.Interval)
	return w
}

// Close shuts down the watcher.
func (w *Watcher) Close() {
	close(w.chanClose)
	<-w.chanClosed
}
