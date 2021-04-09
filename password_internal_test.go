package surveymock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword_Once(t *testing.T) {
	t.Parallel()

	p := newPassword(&Survey{}, "").Once()

	assert.Equal(t, 1, p.repeatability)
}

func TestPassword_Twice(t *testing.T) {
	t.Parallel()

	p := newPassword(&Survey{}, "").Twice()

	assert.Equal(t, 2, p.repeatability)
}

func TestPassword_Times(t *testing.T) {
	t.Parallel()

	p := newPassword(&Survey{}, "").Times(5)

	assert.Equal(t, 5, p.repeatability)
}

func TestPassword_WithHiddenHelp(t *testing.T) {
	t.Parallel()

	p := newPassword(&Survey{}, "original message").WithHiddenHelp("help")

	assert.Equal(t, "original message ", p.message)
	assert.Equal(t, "original message ", p.expectedMessage)
	assert.Equal(t, "help", p.help)
}

func TestPassword_WithHelp(t *testing.T) {
	t.Parallel()

	p := newPassword(&Survey{}, "original message")

	// With Help.
	p.WithHelp("help")

	assert.Equal(t, "original message ", p.message)
	assert.Equal(t, "original message [? for help] ", p.expectedMessage)
	assert.Equal(t, "help", p.help)

	// Clear Help.
	p.WithHelp("")

	assert.Equal(t, "original message ", p.message)
	assert.Equal(t, "original message ", p.expectedMessage)
	assert.Empty(t, p.help)
}
