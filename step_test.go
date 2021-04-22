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

func TestTotalTimes(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario string
		times    []int
		expected int
	}{
		{
			scenario: "nil",
			times:    nil,
			expected: 1,
		},
		{
			scenario: "empty",
			times:    []int{},
			expected: 1,
		},
		{
			scenario: "one element",
			times:    []int{3},
			expected: 3,
		},
		{
			scenario: "multiple elements",
			times:    []int{1, 2, 3},
			expected: 6,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, totalTimes(tc.times...))
		})
	}
}

func TestRepeatSteps(t *testing.T) {
	t.Parallel()

	actual := repeatStep(pressEnter(), 1, 2)
	expected := []Step{pressEnter(), pressEnter(), pressEnter()}

	assert.Equal(t, expected, actual)
}
