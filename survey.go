package surveyexpect

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

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

// Survey is a expectations container and responsible for testing the prompts.
type Survey struct {
	steps Steps

	// test is An optional variable that holds the test struct, to be used for logging and raising error during the
	// tests.
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

// addStep adds a new step to the sequence.
func (s *Survey) addStep(step Step) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.steps.Append(step)
}

// ExpectConfirm expects a ConfirmPrompt.
//
//    Survey.ExpectConfirm("ConfirmPrompt?").
//    	Yes()
func (s *Survey) ExpectConfirm(message string) *ConfirmPrompt {
	e := newConfirm(s, message)

	s.addStep(e)

	return e
}

// ExpectInput expects an InputPrompt.
//
//    Survey.ExpectInput("Enter password:").
//    	Answer("hello world!")
func (s *Survey) ExpectInput(message string) *InputPrompt {
	e := newInput(s, message).Once()

	s.addStep(e)

	return e
}

// ExpectMultiSelect expects a MultiSelectPrompt.
//
//    Survey.ExpectMultiSelect("Enter password:").
//    	Enter()
func (s *Survey) ExpectMultiSelect(message string) *MultiSelectPrompt {
	e := newMultiSelect(s, message)

	s.addStep(e)

	return e
}

// ExpectPassword expects a PasswordPrompt.
//
//    Survey.ExpectPassword("Enter password:").
//    	Answer("hello world!")
func (s *Survey) ExpectPassword(message string) *PasswordPrompt {
	e := newPassword(s, message).Once()

	s.addStep(e)

	return e
}

// ExpectSelect expects a SelectPrompt.
//
//    Survey.ExpectSelect("Enter password:").
//    	Enter()
func (s *Survey) ExpectSelect(message string) *SelectPrompt {
	e := newSelect(s, message)

	s.addStep(e)

	return e
}

// Expect runs an expectation against a given console.
func (s *Survey) Expect(c Console) error {
	if err := s.steps.DoFirst(c); !IsIgnoredError(err) {
		return err
	}

	return nil
}

// answer runs the expectations in background and notifies when it is done.
func (s *Survey) answer(c Console) <-chan struct{} {
	sig := NewSignal()

	go func() {
		defer sig.Notify()

	expectations:
		for {
			select {
			case <-sig.Done():
				// Already closed by timeout.
				break expectations

			default:
				// If not, we run the expectation.
				if err := s.Expect(c); err != nil {
					if !IsNothingTodo(err) {
						s.test.Errorf(err.Error())
					}

					break expectations
				}
			}
		}

		c.ExpectEOF() // nolint: errcheck,gosec
	}()

	// Force close when timeout.
	go func() {
		select {
		case <-time.After(s.timeout):
			s.test.Log("answer timeout exceeded")
			sig.Notify()

		case <-sig.Done():
		}
	}()

	return sig.Done()
}

// ask runs the survey.
func (s *Survey) ask(c Console, fn func(stdio terminal.Stdio)) <-chan struct{} {
	sig := NewSignal()

	go func() {
		defer func() {
			s.test.Log("close console")

			err := c.Tty().Close()
			require.NoError(s.test, err)

			err = c.Close()
			require.NoError(s.test, err)

			sig.Notify()
		}()

		fn(stdio(c))
	}()

	go func() {
		select {
		case <-time.After(s.timeout):
			s.test.Errorf("ask timeout exceeded")
			sig.Notify()

		case <-sig.Done():
			return
		}
	}()

	return sig.Done()
}

// Start starts the survey with a default timeout.
func (s *Survey) Start(fn func(stdio terminal.Stdio)) {
	s.startMu.Lock()
	defer s.startMu.Unlock()

	// Setup a console.
	buf := new(Buffer)
	console, state, err := vt10x.NewVT10XConsole(expect.WithStdout(buf))
	require.NoError(s.test, err)

	// Run the survey in background and close console when it is done.
	askDone := s.ask(console, fn)

	// Run the answer in background.
	// Wait til the survey is done answering.
	<-s.answer(console)
	<-askDone

	s.test.Logf("Raw output: %q\n", buf.String())

	// Dump the terminal's screen.
	s.test.Logf("%s\n", expect.StripTrailingEmptyLines(state.String()))
}

// ExpectationsWereMet checks whether all queued expectations were met in order.
// If any of them was not met - an error is returned.
func (s *Survey) ExpectationsWereMet() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.steps.ExpectationsWereMet()
	if err == nil {
		return nil
	}

	var sb strings.Builder

	sb.WriteString("there are remaining expectations that were not met:\n\n")
	sb.WriteString(err.Error())

	// nolint:goerr113
	return errors.New(sb.String())
}

// ResetExpectations resets all the expectations.
func (s *Survey) ResetExpectations() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.steps.Reset()
}

// stdio returns a terminal.Stdio of the given console.
func stdio(c Console) terminal.Stdio {
	return terminal.Stdio{
		In:  c.Tty(),
		Out: c.Tty(),
		Err: c.Tty(),
	}
}
