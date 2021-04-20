package surveyexpect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequenceExpectation_String(t *testing.T) {
	t.Parallel()

	t.Run("empty expectation", func(t *testing.T) {
		t.Parallel()

		s := sequenceExpectation()

		assert.Empty(t, s.String())
	})

	t.Run("chain expectations", func(t *testing.T) {
		t.Parallel()

		s := sequenceExpectation(
			moveUpAnswer(),
		)

		s.Chain(moveDownAnswer(), enterAnswer(), escAnswer(), tabAnswer(), typeAnswer("hello"))

		expectedResult := `press MOVE UP, press MOVE DOWN, press ENTER, press ESC, press TAB, type "hello"`

		assert.Equal(t, expectedResult, s.String())
	})
}
