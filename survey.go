package surveymock

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/require"
)

// StringWriter is a wrapper for bytes.Buffer.
type StringWriter interface {
	io.Writer
	fmt.Stringer
}

// Survey is a mocked survey.
type Survey struct {
	expectations []Expectation

	// test is An optional variable that holds the test struct, to be used when an
	// invalid mock call was made.
	test TestingT

	timeout time.Duration

	mu      sync.Mutex
	startMu sync.Mutex
}

// WithTimeout sets the timeout of a survey.
func (s *Survey) WithTimeout(t time.Duration) *Survey {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.timeout = t

	return s
}

// expect adds a new expectation to the queue.
func (s *Survey) expect(e Expectation) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.expectations = append(s.expectations, e)
}

// ExpectPassword expects a Password.
//
//    Survey.ExpectPassword("Enter password:").
//    	Answer("hello world!")
func (s *Survey) ExpectPassword(question string) *Password {
	e := newPassword(s, question).Once()

	s.expect(e)

	return e
}

// runExpectation runs an expectation against a given console.
func (s *Survey) runExpectation(c Console) error {
	s.mu.Lock()
	e := s.expectations[0]
	s.mu.Unlock()

	if err := e.Expect(c); err != nil && !errors.Is(err, terminal.InterruptErr) {
		return err
	}

	if e.Repeat() {
		return nil
	}

	// Remove the expectation from the queue if it is not recurrent.
	s.mu.Lock()
	defer s.mu.Unlock()

	s.expectations[0] = nil
	s.expectations = s.expectations[1:]

	return nil
}

// answer runs the expectations in background and notifies when it is done.
func (s *Survey) answer(c Console, rawOutput StringWriter) <-chan struct{} {
	done := make(chan struct{})

	go func() {
	expectations:
		for {
			s.mu.Lock()
			count := len(s.expectations)
			s.mu.Unlock()

			if count == 0 {
				break expectations
			}

			select {
			case <-done:
				// Already closed by timeout.
				break expectations

			default:
				// If not, we run the expectation.
				if err := s.runExpectation(c); err != nil {
					s.test.Errorf(err.Error())

					break expectations
				}
			}
		}

		c.ExpectEOF() // nolint: errcheck,gosec

		s.test.Logf("Raw output: %q\n", rawOutput.String())

		select {
		case <-done:
			// Already closed by timeout.

		default:
			// If not, we close it.
			close(done)
		}
	}()

	// Force close when timeout.
	go func() {
		select {
		case <-time.After(s.timeout):
			s.test.Log("answer timeout exceeded")
			close(done)

		case <-done:
			s.test.Log("answer is done")
		}
	}()

	return done
}

// ask runs the survey.
func (s *Survey) ask(c Console, fn func(stdio terminal.Stdio)) {
	done := make(chan struct{})

	go func() {
		defer func() {
			s.test.Log("close console")

			err := c.Tty().Close()
			require.NoError(s.test, err)

			err = c.Close()
			require.NoError(s.test, err)

			close(done)
		}()

		fn(stdio(c))
	}()

	go func() {
		select {
		case <-time.After(s.timeout):
			s.test.Errorf("ask timeout exceeded")

		case <-done:
			return
		}
	}()
}

// Start starts the survey with a default timeout.
func (s *Survey) Start(fn func(stdio terminal.Stdio)) {
	s.startMu.Lock()
	defer s.startMu.Unlock()

	// Setup a console.
	buf := new(bytes.Buffer)
	console, state, err := vt10x.NewVT10XConsole(expect.WithStdout(buf))
	require.Nil(s.test, err)

	// Run the survey in background and close console when it is done.
	s.ask(console, fn)

	// Run the answer in background.
	// Wait til the survey is done answering.
	<-s.answer(console, buf)

	// Dump the terminal's screen.
	s.test.Logf("%s\n", expect.StripTrailingEmptyLines(state.String()))
}

// ExpectationsWereMet checks whether all queued expectations were met in order.
// If any of them was not met - an error is returned.
func (s *Survey) ExpectationsWereMet() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if len(s.expectations) == 0 {
		return nil
	}

	var sb strings.Builder

	sb.WriteString("there are remaining expectations that were not met:\n")

	for _, e := range s.expectations {
		sb.WriteRune('\n')
		sb.WriteString(e.String())
	}

	// nolint:goerr113
	return errors.New(sb.String())
}

// stdio returns a terminal.Stdio of the given console.
func stdio(c Console) terminal.Stdio {
	return terminal.Stdio{
		In:  c.Tty(),
		Out: c.Tty(),
		Err: c.Tty(),
	}
}

// WithStdio sets survey.WithStdio for a prompt.
func WithStdio(stdio terminal.Stdio) survey.AskOpt {
	return survey.WithStdio(stdio.In, stdio.Out, stdio.Err)
}
