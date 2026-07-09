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

	// ChanAdded receives an Interface when a new interface is added. This
	// cannot be nil and must be a valid channel.
	ChanAdded chan<- Interface

	// ChanRemoved receives an Interface when an interface is removed. This
	// cannot be nil and must be a valid channel.
	ChanRemoved chan<- Interface
}

// Watcher monitors available network interfaces and notifies when one is
// added or removed.
type Watcher struct {
	chanTest    chan any
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
			select {
			case w.chanRemoved <- v:
			case <-w.chanClose:
				return nil
			}
		}
	}
	for k, v := range m2 {
		if _, ok := m[k]; !ok {
			select {
			case w.chanAdded <- v:
			case <-w.chanClose:
				return nil
			}
		}
	}
	return m2
}

func (w *Watcher) run(interval time.Duration) {
	defer close(w.chanClosed)
	defer close(w.chanRemoved)
	defer close(w.chanAdded)
	t := gotime.NewTicker(interval)
	defer t.Stop()
	m := w.diff(map[string]Interface{})
	for {
		select {
		case <-t.C:
			if w.chanTest != nil {
				close(w.chanTest)
			}
			m = w.diff(m)
			if m == nil {
				return
			}
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
