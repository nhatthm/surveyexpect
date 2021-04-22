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
//    Survey.ExpectInput("Enter your name:").
//    	ShowHelp("It's your full name")
func (p *InputPrompt) ShowHelp(help string, options ...string) {
	p.lock()
	defer p.unlock()

	p.answer = helpAnswer(help, options...)
	p.timesLocked(1)
}

// Interrupt marks the answer is interrupted.
//
//    Survey.ExpectInput("Enter your name:").
//    	Interrupt()
func (p *InputPrompt) Interrupt() {
	p.lock()
	defer p.unlock()

	p.answer = interruptAnswer()
	p.timesLocked(1)
}

// Answer sets the answer to the input prompt.
//
//    Survey.ExpectInput("Enter your name:").
//    	Answer("johnny")
func (p *InputPrompt) Answer(answer string) *InputAnswer {
	p.lock()
	defer p.unlock()

	a := newInputAnswer(p, answer)
	p.answer = a

	return a
}

// Type starts a sequence of steps to interact with suggestion mode.
//
//    Survey.ExpectInput("Enter your name:").
//    	Type("johnny")
func (p *InputPrompt) Type(s string) *InputSuggestionSteps {
	p.lock()
	defer p.unlock()

	a := newInputSuggestionSteps(p, typeAnswer(s))
	p.answer = a

	return a
}

// Tab starts a sequence of steps to interact with suggestion mode. Default is 1 when omitted.
//
//    Survey.ExpectInput("Enter your name:").
//    	Tab()
func (p *InputPrompt) Tab(times ...int) *InputSuggestionSteps {
	p.lock()
	defer p.unlock()

	a := newInputSuggestionSteps(p, repeatStep(pressTab(), times...)...)
	p.answer = a

	return a
}

// Do runs the step.
func (p *InputPrompt) Do(c Console) error {
	if _, err := c.ExpectString(p.message); err != nil {
		return err
	}

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
func (p *InputPrompt) String() string {
	var sb stringsBuilder

	sb.WriteLabelLinef("Expect", "Input Prompt").
		WriteLabelLinef("Message", "%q", p.message)

	if steps, ok := p.answer.(*InputSuggestionSteps); ok {
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
//    Survey.ExpectInput("Enter your name:").
//    	Answer("johnny").
//    	Once()
func (p *InputPrompt) Once() *InputPrompt {
	return p.Times(1)
}

// Twice indicates that the message should only be asked twice.
//
//    Survey.ExpectInput("Enter your name:").
//    	Answer("johnny").
//    	Twice()
func (p *InputPrompt) Twice() *InputPrompt {
	return p.Times(2)
}

// Times indicates that the message should only be asked the indicated number of times.
//
//    Survey.ExpectInput("Enter your name:").
//    	Answer("johnny").
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

// InputSuggestionSteps is a sequence of steps when user is in suggestion mode.
type InputSuggestionSteps struct {
	parent *InputPrompt
	steps  *InlineSteps
}

func (a *InputSuggestionSteps) append(steps ...Step) *InputSuggestionSteps {
	a.parent.lock()
	defer a.parent.unlock()

	a.steps.Append(steps...)

	return a
}

// Tab sends the TAB key the indicated times. Default is 1 when omitted.
//
//    Survey.ExpectInput("Enter your name:").
//    	Type("hello").
//    	Tab(5)
func (a *InputSuggestionSteps) Tab(times ...int) *InputSuggestionSteps {
	return a.append(repeatStep(pressTab(), times...)...)
}

// Esc sends the ESC key.
//
//    Survey.ExpectInput("Enter your name:").
//    	Type("hello").
//    	Esc()
func (a *InputSuggestionSteps) Esc() *InputSuggestionSteps {
	return a.append(pressEsc())
}

// Enter sends the ENTER key and ends the sequence.
//
//    Survey.ExpectInput("Enter your name:").
//    	Type("hello").
//    	Enter()
func (a *InputSuggestionSteps) Enter() {
	a.append(pressEnter())
	a.steps.Close()
}

// Interrupt sends ^C and ends the sequence.
//
//    Survey.ExpectInput("Enter your name:").
//    	Type("johnny").
//    	Interrupt()
func (a *InputSuggestionSteps) Interrupt() {
	a.append(pressInterrupt())
	a.steps.Close()
}

// MoveUp sends the ARROW UP key the indicated times. Default is 1 when omitted.
//
//    Survey.ExpectInput("Enter your name:").
//    	Tab().
//    	MoveUp(5)
func (a *InputSuggestionSteps) MoveUp(times ...int) *InputSuggestionSteps {
	return a.append(repeatStep(pressArrowUp(), times...)...)
}

// MoveDown sends the ARROW DOWN key the indicated times. Default is 1 when omitted.
//
//    Survey.ExpectInput("Enter your name:").
//    	Tab().
//    	MoveDown(5)
func (a *InputSuggestionSteps) MoveDown(times ...int) *InputSuggestionSteps {
	return a.append(repeatStep(pressArrowDown(), times...)...)
}

// Delete sends the DELETE key the indicated times. Default is 1 when omitted.
//
//    Survey.ExpectInput("Enter your name:").
//    	Type("johnny").
//    	Delete(5)
func (a *InputSuggestionSteps) Delete(times ...int) *InputSuggestionSteps {
	return a.append(repeatStep(pressDelete(), times...)...)
}

// Type sends a string without enter.
//
//    Survey.ExpectInput("Enter your name:").
//    	Type("johnny").
//    	Tab().
//    	Type(".c").
//    	Enter()
func (a *InputSuggestionSteps) Type(s string) *InputSuggestionSteps {
	return a.append(typeAnswer(s))
}

// ExpectSuggestions expects a list of suggestions.
func (a *InputSuggestionSteps) ExpectSuggestions(suggestions ...string) *InputSuggestionSteps {
	return a.append(expectSelect(suggestions...))
}

// Do runs the step.
func (a *InputSuggestionSteps) Do(c Console) error {
	return a.steps.Do(c)
}

// String represents the answer as a string.
func (a *InputSuggestionSteps) String() string {
	return a.steps.String()
}

func newInputSuggestionSteps(parent *InputPrompt, initialSteps ...Step) *InputSuggestionSteps {
	return &InputSuggestionSteps{
		parent: parent,
		steps:  inlineSteps(initialSteps...),
	}
}
