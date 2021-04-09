package surveymock

// Expectation is an expectation for mocking survey.
type Expectation interface {
	// expect runs the expectation.
	expect(c Console) error

	// repeat tells survey to repeat the same expectation.
	repeat() bool

	// String represents the expectation as a string.
	String() string
}

type base struct {
	parent *Survey

	repeatability int

	// Amount of times this request has been executed.
	totalCalls int // nolint: structcheck
}

func (b *base) lock() {
	b.parent.mu.Lock()
}

func (b *base) unlock() {
	b.parent.mu.Unlock()
}

func (b *base) times(i int) {
	b.repeatability = i
}

func (b *base) timesLocked(i int) {
	b.lock()
	defer b.unlock()

	b.times(i)
}

func (b *base) repeat() bool {
	b.lock()
	defer b.unlock()

	return b.repeatability > 0
}
