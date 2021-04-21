package surveyexpect

import "sync"

// Signal is a safe chan to notify the others.
type Signal struct {
	mu sync.Mutex
	ch chan struct{}
}

// Done checks whether the notification arrives or not.
func (s *Signal) Done() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ch
}

// Notify notifies the listeners.
func (s *Signal) Notify() {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.ch:
	// Do nothing, it was closed.

	default:
		close(s.ch)
	}
}

// NewSignal creates a new Signal.
func NewSignal() *Signal {
	return &Signal{
		ch: make(chan struct{}, 1),
	}
}
