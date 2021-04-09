package surveymock_test

import (
	"fmt"
	"strings"
	"sync"
)

type TestingT struct {
	error strings.Builder
	log   strings.Builder

	clean func()

	mu sync.Mutex
}

func (t *TestingT) Errorf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, _ = fmt.Fprintf(&t.error, format, args...)
}

func (t *TestingT) Log(args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, _ = fmt.Fprintln(&t.log, args...)
}

func (t *TestingT) Logf(format string, args ...interface{}) {
	t.mu.Lock()
	defer t.mu.Unlock()

	_, _ = fmt.Fprintf(&t.log, format, args...)
}

func (t *TestingT) FailNow() {
	panic("failed")
}

func (t *TestingT) Cleanup(clean func()) {
	t.clean = clean
}

func (t *TestingT) ErrorString() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.error.String()
}

func (t *TestingT) LogString() string {
	t.mu.Lock()
	defer t.mu.Unlock()

	return t.log.String()
}

func T() *TestingT {
	return &TestingT{
		clean: func() {},
	}
}
