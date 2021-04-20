package surveyexpect

import (
	"errors"
	"strings"
	"sync"

	"github.com/AlecAivazis/survey/v2/terminal"
)

// Step is an execution step for a survey.
type Step interface {
	// Do runs the step.
	Do(c Console) error

	// String represents the step as a string.
	String() string
}

// Steps is a chain of Step.
type Steps struct {
	steps []Step

	mu sync.Mutex
}

// lock locks the steps from changing its state.
func (s *Steps) lock() {
	s.mu.Lock()
}

// unlock releases the lock.
func (s *Steps) unlock() {
	s.mu.Unlock()
}

// Append appends an expectation to the sequence.
// nolint: unparam
func (s *Steps) Append(more ...Step) *Steps {
	s.lock()
	defer s.unlock()

	s.steps = append(s.steps, more...)

	return s
}

// Do runs the step.
func (s *Steps) Do(c Console) error {
	s.lock()
	cnt := len(s.steps)
	s.unlock()

	if cnt == 0 {
		return ErrNothingToDo
	}

	s.lock()
	step := s.steps[0]
	s.unlock()

	if err := step.Do(c); err != nil {
		isNotFinished := errors.Is(err, ErrNotFinished)
		if !errors.Is(err, terminal.InterruptErr) && !isNotFinished {
			return err
		}

		if isNotFinished {
			return nil
		}
	}

	// Remove the expectation from the queue if it is not recurrent.
	s.lock()
	defer s.unlock()

	s.steps[0] = nil
	s.steps = s.steps[1:]

	return nil
}

// String represents the answer as a string.
func (s *Steps) String() string {
	result := make([]string, 0, len(s.steps))

	for _, s := range s.steps {
		result = append(result, s.String())
	}

	return strings.Join(result, ", ")
}

// Reset removes all the steps.
func (s *Steps) Reset() {
	s.lock()
	defer s.unlock()

	s.steps = nil
}

// ExpectationsWereMet checks whether all queued expectations were met in order.
// If any of them was not met - an error is returned.
func (s *Steps) ExpectationsWereMet() error {
	s.lock()
	defer s.unlock()

	if len(s.steps) == 0 {
		return nil
	}

	var sb strings.Builder

	for _, step := range s.steps {
		sb.WriteRune('\n')
		sb.WriteString(step.String())
	}

	// nolint:goerr113
	return errors.New(sb.String())
}

func steps(steps ...Step) *Steps {
	return &Steps{
		steps: steps,
	}
}
