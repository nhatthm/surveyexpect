package surveymock

// Expectation is an expectation for mocking survey.
type Expectation interface {
	// Expect runs the expectation.
	Expect(c Console) error

	// Repeat tells survey to repeat the same expectation.
	Repeat() bool

	// String represents the expectation as a string.
	String() string
}

type base struct {
	parent *Survey

	repeatability int

	// Amount of times this request has been executed.
	totalCalls int
}

func (b *base) lock() {
	b.parent.mu.Lock()
}

func (b *base) unlock() {
	b.parent.mu.Unlock()
}

func (b *base) times(i int) {
	b.lock()
	defer b.unlock()

	b.timesLocked(i)
}

func (b *base) timesLocked(i int) {
	b.repeatability = i
}

func (b *base) Repeat() bool {
	b.lock()
	defer b.unlock()

	return b.repeatability > 0
}
