package surveymock

import (
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
)

// ReactionTime is to create a small delay to simulate human reaction.
var ReactionTime = 3 * time.Millisecond

// Answer is an expectation for answering a question.
type Answer interface {
	// Expect runs the expectation.
	Expect(c Console) error

	// String represents the answer as a string.
	String() string
}

// NoAnswer sends an empty line to answer the question.
type NoAnswer struct{}

// Expect runs the expectation.
// nolint: errcheck,gosec
func (a *NoAnswer) Expect(c Console) error {
	c.SendLine("")

	return nil
}

// String represents the answer as a string.
func (a *NoAnswer) String() string {
	return "<no answer>"
}

func noAnswer() *NoAnswer {
	return &NoAnswer{}
}

// InterruptAnswer sends an interrupt sequence to terminate the survey.
type InterruptAnswer struct{}

// Expect runs the expectation.
// nolint: errcheck,gosec
func (a *InterruptAnswer) Expect(c Console) error {
	c.SendLine(string(terminal.KeyInterrupt))

	return terminal.InterruptErr
}

// String represents the answer as a string.
func (a *InterruptAnswer) String() string {
	return "<interrupt>"
}

func interruptAnswer() *InterruptAnswer {
	return &InterruptAnswer{}
}

// HelpAnswer sends a ? to show the help.
type HelpAnswer struct {
	help string
}

// Expect runs the expectation.
// nolint: errcheck,gosec
func (a *HelpAnswer) Expect(c Console) error {
	c.SendLine("?")

	if _, err := c.ExpectString(a.help); err != nil {
		return err
	}

	return nil
}

// String represents the answer as a string.
func (a *HelpAnswer) String() string {
	return "?"
}

func helpAnswer(help string) *HelpAnswer {
	return &HelpAnswer{help: help}
}
