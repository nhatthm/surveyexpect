package surveyexpect

import (
	"fmt"
	"strings"
)

var (
	_ Prompt = (*ConfirmPrompt)(nil)
	_ Answer = (*ConfirmAnswer)(nil)
)

// ConfirmPrompt is an expectation of survey.Confirm.
type ConfirmPrompt struct {
	*basePrompt

	message string
	answer  Answer
}

// ShowHelp sets help for the expectation.
//
//    Survey.ExpectConfirm("Are you sure to delete this file?").
//    	ShowHelp("The file will be permanently deleted").
func (c *ConfirmPrompt) ShowHelp(help string, options ...string) {
	c.lock()
	defer c.unlock()

	c.answer = helpAnswer(help, options...)
}

// Interrupt marks the answer is interrupted.
//
//    Survey.ExpectConfirm("Are you sure to delete this file?").
//    	Interrupt().
func (c *ConfirmPrompt) Interrupt() {
	c.lock()
	defer c.unlock()

	c.answer = interruptAnswer()
}

// Yes sets "yes" as the answer to the prompt.
//
//    Survey.ExpectConfirm("Are you sure to delete this file?").
//    	Yes().
func (c *ConfirmPrompt) Yes() {
	c.lock()
	defer c.unlock()

	a := newConfirmAnswer(c, "yes")
	c.answer = a
}

// No sets "no" as the answer to the prompt.
//
//    Survey.ExpectConfirm("Are you sure to delete this file?").
//    	No().
func (c *ConfirmPrompt) No() {
	c.lock()
	defer c.unlock()

	a := newConfirmAnswer(c, "no")
	c.answer = a
}

// Answer sets a custom answer to the prompt.
//
// If the answer is not not empty, the survey expects to have a feedback from the prompt:
//    `Sorry, your reply was invalid: "hello world!" is not a valid answer, please try again.`
//
//    Survey.ExpectConfirm("Are you sure to delete this file?").
//    	Answer("hello world!").
func (c *ConfirmPrompt) Answer(answer string) *ConfirmAnswer {
	c.lock()
	defer c.unlock()

	a := newConfirmAnswer(c, answer)

	if answer != "" {
		a.withFeedback(fmt.Sprintf(`Sorry, your reply was invalid: %q is not a valid answer, please try again.`, answer))
	}

	c.answer = a

	return a
}

// Do runs the step.
func (c *ConfirmPrompt) Do(console Console) error {
	if _, err := console.ExpectString(c.message); err != nil {
		return err
	}

	_ = waitForCursorTwice(console) // nolint: errcheck

	err := c.answer.Do(console)
	if err != nil && !IsInterrupted(err) {
		return err
	}

	c.lock()
	defer c.unlock()

	c.repeatability--
	c.totalCalls++

	return err
}

// String represents the expectation as a string.
func (c *ConfirmPrompt) String() string {
	var sb stringsBuilder

	return sb.WriteLabelLinef("Expect", "Confirm Prompt").
		WriteLabelLinef("Message", "%q", c.message).
		WriteLabelLinef("Answer", c.answer.String()).
		String()
}

// ConfirmAnswer is an answer for confirm question.
type ConfirmAnswer struct {
	parent      *ConfirmPrompt
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

// Do runs the step.
// nolint: errcheck,gosec
func (a *ConfirmAnswer) Do(c Console) error {
	if a.interrupted {
		c.Send(a.answer)
		c.ExpectEOF()

		return nil
	}

	c.SendLine(a.answer)

	if a.feedback != "" {
		if _, err := c.ExpectString(a.feedback); err != nil {
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

func newConfirm(parent *Survey, message string) *ConfirmPrompt {
	return &ConfirmPrompt{
		basePrompt: &basePrompt{
			parent:        parent,
			repeatability: 1,
		},
		message: message,
		answer:  noAnswer(),
	}
}

func newConfirmAnswer(parent *ConfirmPrompt, answer string) *ConfirmAnswer {
	return &ConfirmAnswer{
		parent: parent,
		answer: answer,
	}
}
