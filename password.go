package surveyexpect

import (
	"fmt"
	"strings"
)

var (
	_ Prompt = (*PasswordPrompt)(nil)
	_ Answer = (*PasswordAnswer)(nil)
)

// PasswordPrompt is an expectation of survey.Password.
type PasswordPrompt struct {
	*basePrompt

	message string
	answer  Answer
}

// ShowHelp sets help for the expectation.
//
//    Survey.ExpectPassword("Enter password:").
//    	ShowHelp("Your shiny password").
func (p *PasswordPrompt) ShowHelp(help string) {
	p.lock()
	defer p.unlock()

	p.answer = helpAnswer(help)
	p.timesLocked(1)
}

// Interrupt marks the answer is interrupted.
//
//    Survey.ExpectPassword("Enter password:").
//    	Interrupt().
func (p *PasswordPrompt) Interrupt() {
	p.lock()
	defer p.unlock()

	p.answer = interruptAnswer()
	p.timesLocked(1)
}

// Answer sets the answer to the password prompt.
//
//    Survey.ExpectPassword("Enter password:").
//    	Answer("hello world!").
func (p *PasswordPrompt) Answer(answer string) *PasswordAnswer {
	p.lock()
	defer p.unlock()

	a := newPasswordAnswer(p, answer)
	p.answer = a

	return a
}

// Do runs the step.
func (p *PasswordPrompt) Do(c Console) error {
	_, err := c.ExpectString(p.message)
	if err != nil {
		return err
	}

	_ = waitForCursorTwice(c) // nolint: errcheck

	err = p.answer.Do(c)
	if err != nil && !IsInterrupted(err) {
		return err
	}

	p.lock()
	defer p.unlock()

	p.repeatability--
	p.totalCalls++

	return p.isDoneLocked(err)
}

// String represents the expectation as a string.
func (p *PasswordPrompt) String() string {
	var sb stringsBuilder

	sb.WriteLabelLinef("Expect", "Password Prompt").
		WriteLabelLinef("Message", "%q", p.message).
		WriteLabelLinef("Answer", "%s", p.answer)

	if p.repeatability > 0 && (p.totalCalls != 0 || p.repeatability != 1) {
		sb.WriteLinef("(called: %d time(s), remaining: %d time(s))", p.totalCalls, p.repeatability)
	}

	return sb.String()
}

// Once indicates that the message should only be asked once.
//
//    Survey.ExpectPassword("Enter password:").
//    	Answer("hello world!").
//    	Once()
func (p *PasswordPrompt) Once() *PasswordPrompt {
	return p.Times(1)
}

// Twice indicates that the message should only be asked twice.
//
//    Survey.ExpectPassword("Enter password:").
//    	Answer("hello world!").
//    	Twice()
func (p *PasswordPrompt) Twice() *PasswordPrompt {
	return p.Times(2)
}

// Times indicates that the message should only be asked the indicated number of times.
//
//    Survey.ExpectPassword("Enter password:").
//    	Answer("hello world!").
//    	Times(5)
func (p *PasswordPrompt) Times(i int) *PasswordPrompt {
	p.times(i)

	return p
}

// PasswordAnswer is an answer for password question.
type PasswordAnswer struct {
	parent      *PasswordPrompt
	answer      string
	interrupted bool
}

// Do runs the step.
// nolint: errcheck,gosec
func (a *PasswordAnswer) Do(c Console) error {
	if a.interrupted {
		c.Send(a.answer)
		c.ExpectEOF()

		return nil
	}

	if a.answer == "" {
		c.SendLine("")

		return nil
	}

	c.Send(a.answer)

	// Do asterisks.
	_, err := c.ExpectString(strings.Repeat("*", len(a.answer)))
	if err != nil {
		return err
	}

	c.SendLine("")

	return nil
}

// Interrupted expects the answer will be interrupted.
func (a *PasswordAnswer) Interrupted() {
	a.parent.lock()
	defer a.parent.unlock()

	a.interrupted = true
}

// String represents the answer as a string.
func (a *PasswordAnswer) String() string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "%q", a.answer)

	if a.interrupted {
		_, _ = sb.WriteString(" and get interrupted")
	}

	return sb.String()
}

func newPassword(parent *Survey, message string) *PasswordPrompt {
	return &PasswordPrompt{
		basePrompt: &basePrompt{parent: parent},
		message:    message,
		answer:     noAnswer(),
	}
}

func newPasswordAnswer(parent *PasswordPrompt, answer string) *PasswordAnswer {
	return &PasswordAnswer{
		parent: parent,
		answer: answer,
	}
}
