package surveyexpect

var _ Prompt = (*SelectPrompt)(nil)

// SelectPrompt is an expectation of survey.Select.
type SelectPrompt struct {
	*basePrompt

	message string
	steps   *InlineSteps
}

func (p *SelectPrompt) append(steps ...Step) *SelectPrompt {
	p.lock()
	defer p.unlock()

	p.steps.Append(steps...)

	return p
}

// ShowHelp asks for help and asserts the help text.
//
//    Survey.ExpectSelect("Select a language:").
//    	ShowHelp("Your preferred language")
func (p *SelectPrompt) ShowHelp(help string, options ...string) *SelectPrompt {
	return p.append(pressHelp(help, options...))
}

// Type sends some text to filter the options.
//
//    Survey.ExpectSelect("Select a language:").
//    	Type("Eng")
func (p *SelectPrompt) Type(s string) *SelectPrompt {
	return p.append(typeAnswer(s))
}

// Tab sends the TAB key the indicated times. Default is 1 when omitted.
//
//    Survey.ExpectSelect("Select a language:").
//    	Type("Eng").
//		Tab()
func (p *SelectPrompt) Tab(times ...int) *SelectPrompt {
	return p.append(repeatStep(pressTab(), times...)...)
}

// Interrupt sends ^C and ends the sequence.
//
//    Survey.ExpectSelect("Select a language:").
//		Interrupt()
func (p *SelectPrompt) Interrupt() {
	p.append(pressInterrupt())
	p.steps.Close()
}

// Enter sends the ENTER key and ends the sequence.
//
//    Survey.ExpectSelect("Select a language:").
//    	Type("Eng").
//		Enter()
func (p *SelectPrompt) Enter() {
	p.append(pressEnter())
	p.steps.Close()
}

// Delete sends the DELETE key the indicated times. Default is 1 when omitted.
//
//    Survey.ExpectSelect("Select a language:").
//    	Type("Eng").
//		Delete(3)
func (p *SelectPrompt) Delete(times ...int) *SelectPrompt {
	return p.append(repeatStep(pressDelete(), times...)...)
}

// MoveUp sends the ARROW UP key the indicated times. Default is 1 when omitted.
//
//    Survey.ExpectSelect("Select a language:").
//    	Type("Eng").
//		MoveUp()
func (p *SelectPrompt) MoveUp(times ...int) *SelectPrompt {
	return p.append(repeatStep(pressArrowUp(), times...)...)
}

// MoveDown sends the ARROW DOWN key the indicated times. Default is 1 when omitted.
//
//    Survey.ExpectSelect("Select a language:").
//    	Type("Eng").
//		MoveDown()
func (p *SelectPrompt) MoveDown(times ...int) *SelectPrompt {
	return p.append(repeatStep(pressArrowDown(), times...)...)
}

// ExpectOptions expects a list of options.
//
//    Survey.ExpectSelect("Select a language:").
//    	Type("Eng").
//		ExpectOptions("English")
func (p *SelectPrompt) ExpectOptions(options ...string) *SelectPrompt {
	return p.append(expectSelect(options...))
}

// Do runs the step.
// nolint: errcheck
func (p *SelectPrompt) Do(c Console) error {
	if _, err := c.ExpectString(p.message); err != nil {
		return err
	}

	_, _ = c.ExpectString("\x1b[?25l")

	return p.steps.Do(c)
}

// String represents the expectation as a string.
func (p *SelectPrompt) String() string {
	var sb stringsBuilder

	return sb.WriteLabelLinef("Expect", "Select Prompt").
		WriteLabelLinef("Message", "%q", p.message).
		WriteString(p.steps.String()).
		String()
}

func newSelect(parent *Survey, message string) *SelectPrompt {
	return &SelectPrompt{
		basePrompt: &basePrompt{parent: parent},
		message:    message,
		steps:      inlineSteps(),
	}
}
