package surveyexpect

import (
	"fmt"
	"strings"
)

var (
	_ Prompt = (*MultilinePrompt)(nil)
	_ Answer = (*MultilineAnswer)(nil)
)

// MultilinePrompt is an expectation of survey.Multiline.
type MultilinePrompt struct {
	*basePrompt

	message string
	answer  Step
}

// Interrupt marks the answer is interrupted.
//
//	Survey.ExpectMultiline("Enter your message:").
//		Interrupt()
func (p *MultilinePrompt) Interrupt() {
	p.lock()
	defer p.unlock()

	p.answer = pressInterrupt()
	p.timesLocked(1)
}

// Answer sets the answer to the input prompt.
//
//	Survey.ExpectMultiline("Enter your message:").
//		Answer("hello world")
func (p *MultilinePrompt) Answer(answer string) *MultilineAnswer {
	p.lock()
	defer p.unlock()

	a := newMultilineAnswer(p, answer)
	p.answer = a

	return a
}

// Do runs the step.
func (p *MultilinePrompt) Do(c Console) error {
	if _, err := c.ExpectString(p.message); err != nil {
		return err
	}

	_ = waitForCursorTwice(c) // nolint: errcheck

	err := p.answer.Do(c)
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
func (p *MultilinePrompt) String() string {
	var sb stringsBuilder

	sb.WriteLabelLinef("Expect", "Multiline Prompt").
		WriteLabelLinef("Message", "%q", p.message).
		WriteLabelLinef("Answer", p.answer.String())

	if p.repeatability > 0 && (p.totalCalls != 0 || p.repeatability != 1) {
		sb.WriteLinef("(called: %d time(s), remaining: %d time(s))", p.totalCalls, p.repeatability)
	}

	return sb.String()
}

// Once indicates that the message should only be asked once.
//
//	Survey.ExpectMultiline("Enter your message:").
//		Answer("hello world").
//		Once()
func (p *MultilinePrompt) Once() *MultilinePrompt {
	return p.Times(1)
}

// Twice indicates that the message should only be asked twice.
//
//	Survey.ExpectMultiline("Enter your message:").
//		Answer("hello world").
//		Twice()
func (p *MultilinePrompt) Twice() *MultilinePrompt {
	return p.Times(2)
}

// Times indicates that the message should only be asked the indicated number of times.
//
//	Survey.ExpectMultiline("Enter your message:").
//		Answer("hello world").
//		Times(5)
func (p *MultilinePrompt) Times(i int) *MultilinePrompt {
	p.times(i)

	return p
}

// MultilineAnswer is an answer for password question.
type MultilineAnswer struct {
	parent      *MultilinePrompt
	answer      string
	interrupted bool
}

// Do runs the step.
// nolint: errcheck,gosec
func (a *MultilineAnswer) Do(c Console) error {
	if a.interrupted {
		c.Send(a.answer)
		c.ExpectEOF()

		return nil
	}

	lines := strings.Split(a.answer, "\n")
	lines = append(lines, "")

	if a.answer != "" {
		lines = append(lines, "")
	}

	cnt := len(lines) - 1

	for i, l := range lines {
		c.SendLine(l)

		if i < cnt {
			_ = waitForCursorTwice(c)
		}
	}

	return nil
}

// Interrupted expects the answer will be interrupted.
func (a *MultilineAnswer) Interrupted() {
	a.parent.lock()
	defer a.parent.unlock()

	a.interrupted = true
}

// String represents the answer as a string.
func (a *MultilineAnswer) String() string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "%q", a.answer)

	if a.interrupted {
		_, _ = sb.WriteString(" and get interrupted")
	}

	return sb.String()
}

func newMultiline(parent *Survey, message string) *MultilinePrompt {
	p := &MultilinePrompt{
		basePrompt: &basePrompt{parent: parent},
		message:    message,
	}

	p.answer = newMultilineAnswer(p, "")

	return p
}

func newMultilineAnswer(parent *MultilinePrompt, answer string) *MultilineAnswer {
	return &MultilineAnswer{
		parent: parent,
		answer: answer,
	}
}
