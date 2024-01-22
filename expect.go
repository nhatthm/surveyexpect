package surveyexpect

import (
	"time"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/stretchr/testify/assert"
)

// TestingT is an interface wrapper around *testing.T.
type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Cleanup(f func())
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

// ExpectOption is option for the survey.
type ExpectOption func(s *Survey)

// Expector exp survey.
type Expector func(t TestingT) *Survey

// New creates a new expected survey.
func New(t TestingT, options ...ExpectOption) *Survey {
	s := &Survey{
		test:    t,
		timeout: 3 * time.Second,
	}

	for _, o := range options {
		o(s)
	}

	return s
}

// Expect creates an expected survey with expectations and assures that ExpectationsWereMet() is called.
func Expect(options ...ExpectOption) Expector {
	return func(t TestingT) *Survey {
		// Setup the survey.
		s := New(t, options...)

		t.Cleanup(func() {
			assert.NoError(t, s.ExpectationsWereMet())
		})

		return s
	}
}

func init() { //nolint: gochecknoinits
	// Disable color output for all prompts to simplify testing.
	core.DisableColor = true
}
