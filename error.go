package surveyexpect

import (
	"errors"

	"github.com/AlecAivazis/survey/v2/terminal"
)

var (
	// ErrNothingToDo indicates that there is nothing to do.
	ErrNothingToDo = errors.New("nothing to do")
	// ErrNotFinished indicates that the step is not finished.
	ErrNotFinished = errors.New("step is not finished")
)

// IsIgnoredError checks whether the error is ignored.
func IsIgnoredError(err error) bool {
	if err == nil {
		return true
	}

	return IsInterrupted(err)
}

// IsInterrupted checks if the error is terminal.InterruptErr or not.
func IsInterrupted(err error) bool {
	return errors.Is(err, terminal.InterruptErr)
}

// IsNothingTodo checks if the error is ErrNothingToDo or not.
func IsNothingTodo(err error) bool {
	return errors.Is(err, ErrNothingToDo)
}
