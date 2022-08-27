package surveyexpect

var _ Prompt = (*MultiSelectPrompt)(nil)

// MultiSelectPrompt is an expectation of survey.Select.
type MultiSelectPrompt struct {
	*basePrompt

	message string
	steps   *InlineSteps
}

func (p *MultiSelectPrompt) append(steps ...Step) *MultiSelectPrompt {
	p.lock()
	defer p.unlock()

	p.steps.Append(steps...)

	return p
}

// ShowHelp asks for help and asserts the help text.
//
//	Survey.ExpectMultiSelect("Select a language:").
//		ShowHelp("Your preferred language")
func (p *MultiSelectPrompt) ShowHelp(help string, options ...string) *MultiSelectPrompt {
	return p.append(pressHelp(help, options...))
}

// Type sends some text to filter the options.
//
//	Survey.ExpectMultiSelect("Select a language:").
//		Type("Eng")
func (p *MultiSelectPrompt) Type(s string) *MultiSelectPrompt {
	return p.append(typeAnswer(s))
}

// Tab sends the TAB key the indicated times. Default is 1 when omitted.
//
//	   Survey.ExpectMultiSelect("Select a language:").
//	   	Type("Eng").
//			Tab()
func (p *MultiSelectPrompt) Tab(times ...int) *MultiSelectPrompt {
	return p.append(repeatStep(pressTab(), times...)...)
}

// Interrupt sends ^C and ends the sequence.
//
//	   Survey.ExpectMultiSelect("Select a language:").
//			Interrupt()
func (p *MultiSelectPrompt) Interrupt() {
	p.append(pressInterrupt())
	p.steps.Close()
}

// Enter sends the ENTER key and ends the sequence.
//
//	   Survey.ExpectMultiSelect("Select a language:").
//	   	Type("Eng").
//			Enter()
func (p *MultiSelectPrompt) Enter() {
	p.append(pressEnter())
	p.steps.Close()
}

// Delete sends the DELETE key the indicated times. Default is 1 when omitted.
//
//	   Survey.ExpectMultiSelect("Select a language:").
//	   	Type("Eng").
//			Delete(3)
func (p *MultiSelectPrompt) Delete(times ...int) *MultiSelectPrompt {
	return p.append(repeatStep(pressDelete(), times...)...)
}

// MoveUp sends the ARROW UP key the indicated times. Default is 1 when omitted.
//
//	   Survey.ExpectMultiSelect("Select a language:").
//	   	Type("Eng").
//			MoveUp()
func (p *MultiSelectPrompt) MoveUp(times ...int) *MultiSelectPrompt {
	return p.append(repeatStep(pressArrowUp(), times...)...)
}

// MoveDown sends the ARROW DOWN key the indicated times. Default is 1 when omitted.
//
//	   Survey.ExpectMultiSelect("Select a language:").
//	   	Type("Eng").
//			MoveDown()
func (p *MultiSelectPrompt) MoveDown(times ...int) *MultiSelectPrompt {
	return p.append(repeatStep(pressArrowDown(), times...)...)
}

// Select selects an option. If the option is selected, it will be deselected.
//
//	   Survey.ExpectMultiSelect("Select a language:").
//	   	Type("Eng").
//			Select()
func (p *MultiSelectPrompt) Select() *MultiSelectPrompt {
	return p.append(pressSpace())
}

// SelectNone deselects all filtered options.
//
//	   Survey.ExpectMultiSelect("Select a language:").
//	   	Type("Eng").
//			SelectNone()
func (p *MultiSelectPrompt) SelectNone() *MultiSelectPrompt {
	return p.append(pressArrowLeft())
}

// SelectAll selects all filtered options.
//
//	   Survey.ExpectMultiSelect("Select a language:").
//	   	Type("Eng").
//			SelectAll()
func (p *MultiSelectPrompt) SelectAll() *MultiSelectPrompt {
	return p.append(pressArrowRight())
}

// ExpectOptions expects a list of options.
//
//	   Survey.ExpectMultiSelect("Select a language:").
//	   	Type("Eng").
//			ExpectOptions("English")
func (p *MultiSelectPrompt) ExpectOptions(options ...string) *MultiSelectPrompt {
	return p.append(expectMultiSelect(options...))
}

// Do runs the step.
func (p *MultiSelectPrompt) Do(c Console) error {
	if _, err := c.ExpectString(p.message); err != nil {
		return err
	}

	return p.steps.Do(c)
}

// String represents the expectation as a string.
func (p *MultiSelectPrompt) String() string {
	var sb stringsBuilder

	return sb.WriteLabelLinef("Expect", "MultiSelect Prompt").
		WriteLabelLinef("Message", "%q", p.message).
		WriteString(p.steps.String()).
		String()
}

func newMultiSelect(parent *Survey, message string) *MultiSelectPrompt {
	return &MultiSelectPrompt{
		basePrompt: &basePrompt{parent: parent},
		message:    message,
		steps:      inlineSteps(),
	}
}
