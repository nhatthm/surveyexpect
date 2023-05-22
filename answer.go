package surveyexpect

import (
	"fmt"
	"time"

	"github.com/AlecAivazis/survey/v2/terminal"
)

// ReactionTime is to create a small delay to simulate human reaction.
var ReactionTime = 10 * time.Millisecond

// WaitForReaction creates a small delay to simulate human reaction.
func WaitForReaction() <-chan time.Time {
	return time.After(ReactionTime)
}

// Answer is an expectation for answering a question.
type Answer interface {
	Step
}

// NoAnswer sends an empty line to answer the question.
type NoAnswer struct{}

// Do runs the step.
func (a *NoAnswer) Do(c Console) error {
	c.SendLine("") //nolint: errcheck,gosec

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
func (a *InterruptAnswer) Do(c Console) error {
	c.SendLine(string(terminal.KeyInterrupt)) //nolint: errcheck,gosec

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
func (a *HelpAnswer) Do(c Console) error {
	c.SendLine(a.icon) //nolint: errcheck,gosec

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

// Action sends an action.
type Action struct {
	code   int32
	action string
}

// Do runs the step.
func (a *Action) Do(c Console) error {
	c.Send(string(a.code)) //nolint: errcheck,gosec

	return nil
}

// String represents the answer as a string.
func (a *Action) String() string {
	return fmt.Sprintf("press %s", a.action)
}

func action(code int32, action string) *Action {
	return &Action{
		code:   code,
		action: action,
	}
}

func pressTab() *Action {
	return action(terminal.KeyTab, "TAB")
}

func pressEsc() *Action {
	return action(terminal.KeyEscape, "ESC")
}

func pressEnter() *Action {
	return action(terminal.KeyEnter, "ENTER")
}

func pressArrowUp() *Action {
	return action(terminal.KeyArrowUp, "ARROW UP")
}

func pressArrowDown() *Action {
	return action(terminal.KeyArrowDown, "ARROW DOWN")
}

func pressArrowLeft() *Action {
	return action(terminal.KeyArrowLeft, "ARROW LEFT")
}

func pressArrowRight() *Action {
	return action(terminal.KeyArrowRight, "ARROW RIGHT")
}

func pressInterrupt() *Action {
	return action(terminal.KeyInterrupt, "INTERRUPT")
}

func pressSpace() *Action {
	return action(terminal.KeySpace, "SPACE")
}

func pressDelete() *Action {
	return action(terminal.KeyDelete, "DELETE")
}

// HelpAction sends a ? to show the help.
type HelpAction struct {
	help string
	icon string
}

// Do runs the step.
func (a *HelpAction) Do(c Console) error {
	c.Send(a.icon) //nolint: errcheck,gosec

	if _, err := c.ExpectString(a.help); err != nil {
		return err
	}

	return nil
}

// String represents the answer as a string.
func (a *HelpAction) String() string {
	return fmt.Sprintf("press %q and see %q", a.icon, a.help)
}

func pressHelp(help string, options ...string) *HelpAction {
	if len(options) == 0 {
		options = append(options, "?")
	}

	return &HelpAction{
		help: help,
		icon: options[0],
	}
}

// TypeAnswer types an answer.
type TypeAnswer struct {
	answer string
}

// Do runs the step.
func (a *TypeAnswer) Do(c Console) error {
	c.Send(a.answer) //nolint: errcheck,gosec

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
