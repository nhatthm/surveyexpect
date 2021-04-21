package surveyexpect

import (
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
)

// ReactionTime is to create a small delay to simulate human reaction.
var ReactionTime = 10 * time.Millisecond

// Answer is an expectation for answering a question.
type Answer interface {
	Step
}

// NoAnswer sends an empty line to answer the question.
type NoAnswer struct{}

// Do runs the step.
// nolint: errcheck,gosec
func (a *NoAnswer) Do(c Console) error {
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

// Do runs the step.
// nolint: errcheck,gosec
func (a *InterruptAnswer) Do(c Console) error {
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
	icon string
}

// Do runs the step.
// nolint: errcheck,gosec
func (a *HelpAnswer) Do(c Console) error {
	c.SendLine(a.icon)

	if _, err := c.ExpectString(a.help); err != nil {
		return err
	}

	return nil
}

// String represents the answer as a string.
func (a *HelpAnswer) String() string {
	return a.icon
}

func helpAnswer(help string, options ...string) *HelpAnswer {
	if len(options) == 0 {
		options = append(options, "?")
	}

	return &HelpAnswer{
		help: help,
		icon: options[0],
	}
}

// ActionAnswer sends an action.
type ActionAnswer struct {
	code   int32
	action string
}

// Do runs the step.
// nolint: errcheck,gosec
func (a *ActionAnswer) Do(c Console) error {
	c.Send(string(a.code))

	return nil
}

// String represents the answer as a string.
func (a *ActionAnswer) String() string {
	return a.action
}

func actionAnswer(code int32, action string) *ActionAnswer {
	return &ActionAnswer{
		code:   code,
		action: action,
	}
}

func tabAnswer() *ActionAnswer {
	return actionAnswer(terminal.KeyTab, "press TAB")
}

func escAnswer() *ActionAnswer {
	return actionAnswer(terminal.KeyEscape, "press ESC")
}

func enterAnswer() *ActionAnswer {
	return actionAnswer(terminal.KeyEnter, "press ENTER")
}

func moveUpAnswer() *ActionAnswer {
	return actionAnswer(terminal.KeyEnter, "press MOVE UP")
}

func moveDownAnswer() *ActionAnswer {
	return actionAnswer(terminal.KeyEnter, "press MOVE DOWN")
}

// TypeAnswer types an answer.
type TypeAnswer struct {
	answer string
}

// Do runs the step.
// nolint: errcheck,gosec
func (a *TypeAnswer) Do(c Console) error {
	c.Send(a.answer)

	return nil
}

// String represents the answer as a string.
func (a *TypeAnswer) String() string {
	return fmt.Sprintf("type %q", a.answer)
}

func typeAnswer(answer string) *TypeAnswer {
	return &TypeAnswer{
		answer: answer,
	}
}
