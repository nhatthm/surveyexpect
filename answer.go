package surveymock

import "github.com/AlecAivazis/survey/v2/terminal"

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
	c.Send(string(terminal.KeyInterrupt))
	c.SendLine("")

	return terminal.InterruptErr
}

// String represents the answer as a string.
func (a *InterruptAnswer) String() string {
	return "<interrupt>"
}

func interruptAnswer() *InterruptAnswer {
	return &InterruptAnswer{}
}
