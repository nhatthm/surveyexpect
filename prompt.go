package surveyexpect

// Prompt is a prompt expectation for a survey.
type Prompt interface {
	Step
}

type basePrompt struct {
	parent *Survey

	repeatability int

	// Amount of times this request has been executed.
	totalCalls int
}

func (p *basePrompt) lock() {
	p.parent.mu.Lock()
}

func (p *basePrompt) unlock() {
	p.parent.mu.Unlock()
}

func (p *basePrompt) times(i int) {
	p.lock()
	defer p.unlock()

	p.timesLocked(i)
}

func (p *basePrompt) timesLocked(i int) {
	p.repeatability = i
}

func (p *basePrompt) isDoneLocked(err error) error {
	if err != nil {
		return err
	}

	if p.repeatability > 0 {
		return ErrNotFinished
	}

	return nil
}
