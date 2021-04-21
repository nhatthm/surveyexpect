package surveyexpect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSteps_Append(t *testing.T) {
	t.Parallel()

	t.Run("still open", func(t *testing.T) {
		t.Parallel()

		s := steps()

		assert.NotPanics(t, func() {
			s.Append(pressTab())
		})
	})

	t.Run("closed", func(t *testing.T) {
		t.Parallel()

		s := steps()
		s.Close()

		assert.Panics(t, func() {
			s.Append(pressTab())
		})
	})
}

func TestSteps_String(t *testing.T) {
	t.Parallel()

	t.Run("empty steps", func(t *testing.T) {
		t.Parallel()

		s := steps()

		assert.Empty(t, s.String())
	})

	t.Run("append empty steps", func(t *testing.T) {
		t.Parallel()

		s := steps()

		s.Append(pressArrowDown())

		expectedResult := "press ARROW DOWN"

		assert.Equal(t, expectedResult, s.String())
	})

	t.Run("append existing steps", func(t *testing.T) {
		t.Parallel()

		s := steps(
			pressArrowUp(),
		)

		s.Append(pressArrowDown(), pressEnter(), pressEsc(), pressTab(), typeAnswer("hello"))

		expectedResult := "press ARROW UP\n\npress ARROW DOWN\n\npress ENTER\n\npress ESC\n\npress TAB\n\ntype \"hello\""

		assert.Equal(t, expectedResult, s.String())
	})

	t.Run("reset and re-append", func(t *testing.T) {
		t.Parallel()

		s := steps(pressArrowUp())
		s.Reset()
		s.Append(pressArrowDown())

		expectedResult := "press ARROW DOWN"

		assert.Equal(t, expectedResult, s.String())
	})
}
