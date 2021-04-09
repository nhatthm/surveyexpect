package surveymock

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2/terminal"
)

var _ Expectation = (*Password)(nil)

// Password is an expectation of survey.Password.
type Password struct {
	*base

	message         string
	expectedMessage string
	help            string
	answer          string
}

// WithHiddenHelp sets help for the expectation.
//
//    Survey.ExpectPassword("Enter password:").
//    	WithHiddenHelp("Your shiny password").
//    	Answer("hello world!").
func (p *Password) WithHiddenHelp(help string) *Password {
	p.lock()
	defer p.unlock()

	p.help = help

	return p
}

// WithHelp sets help for the expectation.
//
//    Survey.ExpectPassword("Enter password:").
//    	WithHelp("Your shiny password").
//    	Answer("hello world!").
func (p *Password) WithHelp(help string) *Password {
	p.lock()
	defer p.unlock()

	p.help = help

	if help == "" {
		p.expectedMessage = p.message
	} else {
		p.expectedMessage = fmt.Sprintf("%s[? for help] ", p.message)
	}

	return p
}

// Interrupt marks the answer is interrupted.
//
//    Survey.ExpectPassword("Enter password:").
//    	Interrupt().
func (p *Password) Interrupt() *Password {
	p.lock()
	defer p.unlock()

	p.answer = string(terminal.KeyInterrupt)
	p.times(1)

	return p
}

// Answer sets the answer to the password prompt.
//
//    Survey.ExpectPassword("Enter password:").
//    	Answer("hello world!").
func (p *Password) Answer(answer string) *Password {
	p.lock()
	defer p.unlock()

	p.answer = answer

	return p
}

// expect runs the expectation.
// nolint: errcheck,gosec
func (p *Password) expect(c Console) error {
	_, err := c.ExpectString(p.expectedMessage)
	if err != nil {
		return err
	}

	if p.help != "" {
		c.SendLine("?")

		_, err := c.ExpectString(p.help)
		if err != nil {
			return err
		}
	}

	if p.answer != "" {
		c.Send(p.answer)
	}

	c.SendLine("")

	p.repeatability--
	p.totalCalls++

	return nil
}

// String represents the expectation as a string.
func (p *Password) String() string {
	var sb strings.Builder

	_, _ = sb.WriteString("Type   : Password\n")
	_, _ = sb.WriteString("Message: ")
	_, _ = sb.WriteString(p.expectedMessage)
	_, _ = sb.WriteRune('\n')

	if p.help != "" {
		_, _ = sb.WriteString("Help   : ")
		_, _ = sb.WriteString(p.help)
		_, _ = sb.WriteRune('\n')
	}

	if p.answer != "" {
		_, _ = sb.WriteString("Answer : ")
		_, _ = sb.WriteString(p.answer)
		_, _ = sb.WriteRune('\n')
	}

	if p.repeatability > 0 && (p.totalCalls != 0 || p.repeatability != 1) {
		_, _ = fmt.Fprintf(&sb, "(called: %d time(s), remaining: %d time(s))", p.totalCalls, p.repeatability)
		_, _ = sb.WriteRune('\n')
	}

	return sb.String()
}

func newPassword(parent *Survey, message string) *Password {
	message += " "

	return &Password{
		base:            &base{parent: parent},
		message:         message,
		expectedMessage: message,
	}
}

// Once indicates that the message should only be asked once.
//
//    Survey.ExpectPassword("Enter password:").
//    	Answer("hello world!").
//    	Once()
func (p *Password) Once() *Password {
	return p.Times(1)
}

// Twice indicates that the message should only be asked twice.
//
//    Survey.ExpectPassword("Enter password:").
//    	Answer("hello world!").
//    	Twice()
func (p *Password) Twice() *Password {
	return p.Times(2)
}

// Times indicates that the message should only be asked the indicated number of times.
//
//    Survey.ExpectPassword("Enter password:").
//    	Answer("hello world!").
//    	Times(5)
func (p *Password) Times(i int) *Password {
	p.timesLocked(i)

	return p
}
