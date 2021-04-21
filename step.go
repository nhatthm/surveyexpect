package surveyexpect

import (
	"errors"
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

// DoFirst runs the first step.
func (s *Steps) DoFirst(c Console) error {
	if s.HasNothingToDo() {
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

// Do runs all the steps.
func (s *Steps) Do(c Console) error {
	for {
		if err := s.DoFirst(c); err != nil {
			if IsNothingTodo(err) {
				return nil
			}

			if !IsInterrupted(err) {
				return err
			}
		}
	}
}

// String represents the answer as a string.
func (s *Steps) String() string {
	if s.HasNothingToDo() {
		return ""
	}

	s.lock()
	defer s.unlock()

	var sb stringsBuilder

	for _, step := range s.steps {
		sb.WriteRune('\n').
			WriteString(step.String())
	}

	return sb.String()
}

// Reset removes all the steps.
func (s *Steps) Reset() {
	s.lock()
	defer s.unlock()

	s.steps = nil
}

// Len returns the number of steps.
func (s *Steps) Len() int {
	s.lock()
	defer s.unlock()

	return len(s.steps)
}

// HasNothingToDo checks whether the steps are not empty.
func (s *Steps) HasNothingToDo() bool {
	return s.Len() == 0
}

// ExpectationsWereMet checks whether all queued expectations were met in order.
// If any of them was not met - an error is returned.
func (s *Steps) ExpectationsWereMet() error {
	if s.HasNothingToDo() {
		return nil
	}

	// nolint:goerr113
	return errors.New(s.String())
}

func steps(steps ...Step) *Steps {
	return &Steps{
		steps: steps,
	}
}
