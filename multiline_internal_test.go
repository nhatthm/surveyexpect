package surveyexpect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultilinePrompt_Once(t *testing.T) {
	t.Parallel()

	p := newMultiline(&Survey{}, "").Once()

	assert.Equal(t, 1, p.repeatability)
}

func TestMultilinePrompt_Twice(t *testing.T) {
	t.Parallel()

	p := newMultiline(&Survey{}, "").Twice()

	assert.Equal(t, 2, p.repeatability)
}

func TestMultilinePrompt_Times(t *testing.T) {
	t.Parallel()

	p := newMultiline(&Survey{}, "").Times(5)

	assert.Equal(t, 5, p.repeatability)
}

func TestMultilinePrompt_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		repeatability int
		totalCalls    int
		expected      string
	}{
		{
			scenario: "repeat = 0",
			expected: "Expect : Multiline Prompt\nMessage: \"Enter the password:\"\nAnswer : <no answer>\n",
		},
		{
			scenario:      "repeat == 1 and called = 0",
			repeatability: 1,
			expected:      "Expect : Multiline Prompt\nMessage: \"Enter the password:\"\nAnswer : <no answer>\n",
		},
		{
			scenario:      "repeat > 0",
			repeatability: 3,
			totalCalls:    1,
			expected:      "Expect : Multiline Prompt\nMessage: \"Enter the password:\"\nAnswer : <no answer>\n(called: 1 time(s), remaining: 3 time(s))\n",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			p := &MultilinePrompt{
				basePrompt: &basePrompt{
					repeatability: tc.repeatability,
					totalCalls:    tc.totalCalls,
				},
				message: "Enter the password:",
				answer:  noAnswer(),
			}

			assert.Equal(t, tc.expected, p.String())
		})
	}
}

func TestMultilineAnswer_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario    string
		interrupted bool
		expected    string
	}{
		{
			scenario: "not interrupted",
			expected: `"password"`,
		},
		{
			scenario:    "interrupted",
			interrupted: true,
			expected:    `"password" and get interrupted`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			a := &MultilineAnswer{
				answer:      "password",
				interrupted: tc.interrupted,
			}

			assert.Equal(t, tc.expected, a.String())
		})
	}
}
