package surveyexpect

import (
	"fmt"
	"strings"
)

var (
	_ Prompt = (*InputPrompt)(nil)
	_ Answer = (*InputAnswer)(nil)
)

// InputPrompt is an expectation of survey.Input.
type InputPrompt struct {
	*basePrompt

	message string
	answer  Step
}

// ShowHelp sets help for the expectation.
//
//    Survey.ExpectInput("Enter password:").
//    	ShowHelp("Your shiny password").
func (p *InputPrompt) ShowHelp(help string, options ...string) {
	p.lock()
	defer p.unlock()

	p.answer = helpAnswer(help, options...)
	p.timesLocked(1)
}

// Interrupt marks the answer is interrupted.
//
//    Survey.ExpectInput("Enter password:").
//    	Interrupt().
func (p *InputPrompt) Interrupt() {
	p.lock()
	defer p.unlock()

	p.answer = interruptAnswer()
	p.timesLocked(1)
}

// Answer sets the answer to the password prompt.
//
//    Survey.ExpectInput("Enter password:").
//    	Answer("hello world!").
func (p *InputPrompt) Answer(answer string) *InputAnswer {
	p.lock()
	defer p.unlock()

	a := newInputAnswer(p, answer)
	p.answer = a

	return a
}

// Do runs the step.
func (p *InputPrompt) Do(c Console) error {
	_, err := c.ExpectString(p.message)
	if err != nil {
		return err
	}

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
func (p *InputPrompt) String() string {
	var sb stringsBuilder

	sb.WriteLabelLinef("Expect", "Input Prompt").
		WriteLabelLinef("Message", "%q", p.message)

	if steps, ok := p.answer.(*Steps); ok {
		sb.WriteString(steps.String())
	} else {
		sb.WriteLabelLinef("Answer", p.answer.String())
	}

	if p.repeatability > 0 && (p.totalCalls != 0 || p.repeatability != 1) {
		sb.WriteLinef("(called: %d time(s), remaining: %d time(s))", p.totalCalls, p.repeatability)
	}

	return sb.String()
}

// Once indicates that the message should only be asked once.
//
//    Survey.ExpectInput("Enter password:").
//    	Answer("hello world!").
//    	Once()
func (p *InputPrompt) Once() *InputPrompt {
	return p.Times(1)
}

// Twice indicates that the message should only be asked twice.
//
//    Survey.ExpectInput("Enter password:").
//    	Answer("hello world!").
//    	Twice()
func (p *InputPrompt) Twice() *InputPrompt {
	return p.Times(2)
}

// Times indicates that the message should only be asked the indicated number of times.
//
//    Survey.ExpectInput("Enter password:").
//    	Answer("hello world!").
//    	Times(5)
func (p *InputPrompt) Times(i int) *InputPrompt {
	p.times(i)

	return p
}

// InputAnswer is an answer for password question.
type InputAnswer struct {
	parent      *InputPrompt
	answer      string
	interrupted bool
}

// Do runs the step.
// nolint: errcheck,gosec
func (a *InputAnswer) Do(c Console) error {
	if a.interrupted {
		c.Send(a.answer)
		c.ExpectEOF()

		return nil
	}

	c.SendLine(a.answer)

	return nil
}

// Interrupted expects the answer will be interrupted.
func (a *InputAnswer) Interrupted() {
	a.parent.lock()
	defer a.parent.unlock()

	a.interrupted = true
}

// String represents the answer as a string.
func (a *InputAnswer) String() string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "%q", a.answer)

	if a.interrupted {
		_, _ = sb.WriteString(" and get interrupted")
	}

	return sb.String()
}

func newInput(parent *Survey, message string) *InputPrompt {
	return &InputPrompt{
		basePrompt: &basePrompt{parent: parent},
		message:    message,
		answer:     noAnswer(),
	}
}

func newInputAnswer(parent *InputPrompt, answer string) *InputAnswer {
	return &InputAnswer{
		parent: parent,
		answer: answer,
	}
}
