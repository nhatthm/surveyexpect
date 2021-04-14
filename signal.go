package surveyexpect

import "sync"

// signalCh is a safe chan to notify the others.
type signalCh struct {
	mu sync.Mutex
	ch chan struct{}
}

func (s *signalCh) done() <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.ch
}

func (s *signalCh) close() {
	s.mu.Lock()
	defer s.mu.Unlock()

	select {
	case <-s.ch:
	// Do nothing, it was closed.

	default:
		close(s.ch)
	}
}

func signal() *signalCh {
	return &signalCh{
		ch: make(chan struct{}, 1),
	}
}
