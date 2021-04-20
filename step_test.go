package surveyexpect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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

		expectedResult := `press MOVE DOWN`

		assert.Equal(t, expectedResult, s.String())
	})

	t.Run("append existing steps", func(t *testing.T) {
		t.Parallel()

		s := steps(
			moveUpAnswer(),
		)

		s.Append(moveDownAnswer(), enterAnswer(), escAnswer(), tabAnswer(), typeAnswer("hello"))

		expectedResult := `press MOVE UP, press MOVE DOWN, press ENTER, press ESC, press TAB, type "hello"`

		assert.Equal(t, expectedResult, s.String())
	})

	t.Run("reset and re-append", func(t *testing.T) {
		t.Parallel()

		s := steps(moveUpAnswer())
		s.Reset()
		s.Append(moveDownAnswer())

		expectedResult := `press MOVE DOWN`

		assert.Equal(t, expectedResult, s.String())
	})
}
