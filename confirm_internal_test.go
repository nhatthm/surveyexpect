package surveyexpect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfirm_String(t *testing.T) {
	t.Parallel()

	expected := "Type   : ConfirmPrompt\nMessage: \"ConfirmPrompt?\"\nAnswer : <no answer>\n"

	c := &ConfirmPrompt{
		message: "ConfirmPrompt?",
		answer:  noAnswer(),
	}

	assert.Equal(t, expected, c.String())
}

func TestConfirmAnswer_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario    string
		interrupted bool
		feedback    string
		expected    string
	}{
		{
			scenario:    "interrupted",
			interrupted: true,
			expected:    "\"answer\" and get interrupted",
		},
		{
			scenario: "feedback",
			feedback: "feedback",
			expected: "\"answer\" and get feedback \"feedback\"",
		},
		{
			scenario: "normal",
			expected: `"answer"`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			a := &ConfirmAnswer{
				answer:      "answer",
				feedback:    tc.feedback,
				interrupted: tc.interrupted,
			}

			assert.Equal(t, tc.expected, a.String())
		})
	}
}
