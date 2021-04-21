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
			s.Append(tabAnswer())
		})
	})

	t.Run("closed", func(t *testing.T) {
		t.Parallel()

		s := steps()
		s.Close()

		assert.Panics(t, func() {
			s.Append(tabAnswer())
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

		s.Append(moveDownAnswer())

		expectedResult := "press MOVE DOWN"

		assert.Equal(t, expectedResult, s.String())
	})

	t.Run("append existing steps", func(t *testing.T) {
		t.Parallel()

		s := steps(
			moveUpAnswer(),
		)

		s.Append(moveDownAnswer(), enterAnswer(), escAnswer(), tabAnswer(), typeAnswer("hello"))

		expectedResult := "press MOVE UP\n\npress MOVE DOWN\n\npress ENTER\n\npress ESC\n\npress TAB\n\ntype \"hello\""

		assert.Equal(t, expectedResult, s.String())
	})

	t.Run("reset and re-append", func(t *testing.T) {
		t.Parallel()

		s := steps(moveUpAnswer())
		s.Reset()
		s.Append(moveDownAnswer())

		expectedResult := "press MOVE DOWN"

		assert.Equal(t, expectedResult, s.String())
	})
}
