package surveymock

import (
	"time"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/stretchr/testify/assert"
)

// TestingT is an interface wrapper around *testing.T.
type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Cleanup(func())
	Log(args ...interface{})
	Logf(format string, args ...interface{})
}

// MockOption is option for mocking survey.
type MockOption func(s *Survey)

// Mocker mocks survey.
type Mocker func(t TestingT) *Survey

// New creates a new mocked survey.
func New(t TestingT, mocks ...MockOption) *Survey {
	s := &Survey{
		test:    t,
		timeout: 3 * time.Second,
	}

	for _, m := range mocks {
		m(s)
	}

	return s
}

// Mock creates a mocked server with expectations and assures that ExpectationsWereMet() is called.
func Mock(mocks ...MockOption) Mocker {
	return func(t TestingT) *Survey {
		// Setup mocked survey.
		s := New(t, mocks...)

		t.Cleanup(func() {
			assert.NoError(t, s.ExpectationsWereMet())
		})

		return s
	}
}

// nolint: gochecknoinits
func init() {
	// Disable color output for all prompts to simplify testing.
	core.DisableColor = true
}
