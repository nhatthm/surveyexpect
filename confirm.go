package surveymock

import (
	"errors"
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
)

var (
	_ Expectation = (*Confirm)(nil)
	_ Answer      = (*ConfirmAnswer)(nil)
)

// Confirm is an expectation of survey.Confirm.
type Confirm struct {
	*base

	message string
	answer  Answer
}

// ShowHelp sets help for the expectation.
//
//    Survey.ExpectConfirm("Are you sure to delete this file?").
//    	ShowHelp("The file will be permanently deleted").
func (c *Confirm) ShowHelp(help string) {
	c.lock()
	defer c.unlock()

	c.answer = helpAnswer(help)
}

// Interrupt marks the answer is interrupted.
//
//    Survey.ExpectConfirm("Are you sure to delete this file?").
//    	Interrupt().
func (c *Confirm) Interrupt() {
	c.lock()
	defer c.unlock()

	c.answer = interruptAnswer()
}

// Yes sets "yes" as the answer to the prompt.
//
//    Survey.ExpectConfirm("Are you sure to delete this file?").
//    	Yes().
func (c *Confirm) Yes() {
	c.lock()
	defer c.unlock()

	a := newConfirmAnswer(c, "yes")
	c.answer = a
}

// No sets "no" as the answer to the prompt.
//
//    Survey.ExpectConfirm("Are you sure to delete this file?").
//    	No().
func (c *Confirm) No() {
	c.lock()
	defer c.unlock()

	a := newConfirmAnswer(c, "no")
	c.answer = a
}

// Answer sets a custom answer to the prompt.
//
// If the answer is not not empty, the mock expects to have a feedback from the survey:
//    `Sorry, your reply was invalid: "hello world!" is not a valid answer, please try again.`
//
//    Survey.ExpectConfirm("Are you sure to delete this file?").
//    	Answer("hello world!").
func (c *Confirm) Answer(answer string) *ConfirmAnswer {
	c.lock()
	defer c.unlock()

	a := newConfirmAnswer(c, answer)

	if answer != "" {
		a.withFeedback(fmt.Sprintf(`Sorry, your reply was invalid: %q is not a valid answer, please try again.`, answer))
	}

	c.answer = a

	return a
}

// Expect runs the expectation.
func (c *Confirm) Expect(console Console) error {
	_, err := console.ExpectString(c.message)
	if err != nil {
		return err
	}

	_ = waitForCursorTwice(console) // nolint: errcheck

	err = c.answer.Expect(console)
	if err != nil && !errors.Is(err, terminal.InterruptErr) {
		return err
	}

	c.lock()
	defer c.unlock()

	c.repeatability--
	c.totalCalls++

	return err
}

// String represents the expectation as a string.
func (c *Confirm) String() string {
	var sb strings.Builder

	_, _ = sb.WriteString("Type   : Confirm\n")
	_, _ = fmt.Fprintf(&sb, "Message: %q\n", c.message)
	_, _ = fmt.Fprintf(&sb, "Answer : %s\n", c.answer.String())

	return sb.String()
}

// ConfirmAnswer is an answer for confirm question.
type ConfirmAnswer struct {
	parent      *Confirm
	answer      string
	feedback    string
	interrupted bool
}

func (a *ConfirmAnswer) withFeedback(feedback string) *ConfirmAnswer {
	a.feedback = feedback
	a.interrupted = false

	return a
}

// Interrupted expects the answer will be interrupted.
func (a *ConfirmAnswer) Interrupted() {
	a.parent.lock()
	defer a.parent.unlock()

	a.interrupted = true
	a.feedback = ""
}

// Expect runs the expectation.
// nolint: errcheck,gosec
func (a *ConfirmAnswer) Expect(c Console) error {
	if a.interrupted {
		c.Send(a.answer)
		c.ExpectEOF()

		return nil
	}

	c.SendLine(a.answer)

	if a.feedback != "" {
		_, err := c.ExpectString(a.feedback)
		if err != nil {
			return err
		}
	}

	return nil
}

// String represents the answer as a string.
func (a *ConfirmAnswer) String() string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "%q", a.answer)

	if a.interrupted {
		_, _ = sb.WriteString(" and get interrupted")
	} else if a.feedback != "" {
		_, _ = fmt.Fprintf(&sb, " and get feedback %q", a.feedback)
	}

	return sb.String()
}

func newConfirm(parent *Survey, message string) *Confirm {
	return &Confirm{
		base: &base{
			parent:        parent,
			repeatability: 1,
		},
		message: message,
		answer:  noAnswer(),
	}
}

func newConfirmAnswer(parent *Confirm, answer string) *ConfirmAnswer {
	return &ConfirmAnswer{
		parent: parent,
		answer: answer,
	}
}
