package surveyexpect

import "strings"

// Expectation is an expectation for a survey.
type Expectation interface {
	// Expect runs the expectation.
	Expect(c Console) error

	// String represents the expectation as a string.
	String() string
}

// SequenceExpectation is a chain of answers.
type SequenceExpectation struct {
	sequences []Expectation
}

// Chain appends an expectation to the sequence.
func (a *SequenceExpectation) Chain(more ...Expectation) *SequenceExpectation {
	a.sequences = append(a.sequences, more...)

	return a
}

// Expect runs the expectation.
// nolint: errcheck
func (a *SequenceExpectation) Expect(c Console) error {
	for _, s := range a.sequences {
		if err := s.Expect(c); err != nil {
			return err
		}

		_ = waitForCursorTwice(c)
	}

	return nil
}

// String represents the answer as a string.
func (a *SequenceExpectation) String() string {
	result := make([]string, 0, len(a.sequences))

	for _, s := range a.sequences {
		result = append(result, s.String())
	}

	return strings.Join(result, ", ")
}

func sequenceExpectation(chain ...Expectation) *SequenceExpectation {
	return &SequenceExpectation{
		sequences: chain,
	}
}
