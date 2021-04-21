package surveyexpect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInputPrompt_Once(t *testing.T) {
	t.Parallel()

	p := newInput(&Survey{}, "").Once()

	assert.Equal(t, 1, p.repeatability)
}

func TestInputPrompt_Twice(t *testing.T) {
	t.Parallel()

	p := newInput(&Survey{}, "").Twice()

	assert.Equal(t, 2, p.repeatability)
}

func TestInputPrompt_Times(t *testing.T) {
	t.Parallel()

	p := newInput(&Survey{}, "").Times(5)

	assert.Equal(t, 5, p.repeatability)
}

func TestInputPrompt_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		repeatability int
		totalCalls    int
		expected      string
	}{
		{
			scenario: "repeat = 0",
			expected: "Expect : Input Prompt\nMessage: \"Enter the username:\"\nAnswer : <no answer>\n",
		},
		{
			scenario:      "repeat == 1 and called = 0",
			repeatability: 1,
			expected:      "Expect : Input Prompt\nMessage: \"Enter the username:\"\nAnswer : <no answer>\n",
		},
		{
			scenario:      "repeat > 0",
			repeatability: 3,
			totalCalls:    1,
			expected:      "Expect : Input Prompt\nMessage: \"Enter the username:\"\nAnswer : <no answer>\n(called: 1 time(s), remaining: 3 time(s))\n",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			p := &InputPrompt{
				basePrompt: &basePrompt{
					repeatability: tc.repeatability,
					totalCalls:    tc.totalCalls,
				},
				message: "Enter the username:",
				answer:  noAnswer(),
			}

			assert.Equal(t, tc.expected, p.String())
		})
	}
}

func TestInputAnswer_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario    string
		interrupted bool
		expected    string
	}{
		{
			scenario: "not interrupted",
			expected: `"username"`,
		},
		{
			scenario:    "interrupted",
			interrupted: true,
			expected:    `"username" and get interrupted`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			a := &InputAnswer{
				answer:      "username",
				interrupted: tc.interrupted,
			}

			assert.Equal(t, tc.expected, a.String())
		})
	}
}
