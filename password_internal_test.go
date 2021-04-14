package surveyexpect

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

func TestPassword_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		repeatability int
		totalCalls    int
		expected      string
	}{
		{
			scenario: "repeat = 0",
			expected: "Type   : Password\nMessage: \"Enter the password:\"\nAnswer : <no answer>\n",
		},
		{
			scenario:      "repeat == 1 and called = 0",
			repeatability: 1,
			expected:      "Type   : Password\nMessage: \"Enter the password:\"\nAnswer : <no answer>\n",
		},
		{
			scenario:      "repeat > 0",
			repeatability: 3,
			totalCalls:    1,
			expected:      "Type   : Password\nMessage: \"Enter the password:\"\nAnswer : <no answer>\n(called: 1 time(s), remaining: 3 time(s))\n",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			p := &Password{
				base: &base{
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

func TestPasswordAnswer_String(t *testing.T) {
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

			a := &PasswordAnswer{
				answer:      "password",
				interrupted: tc.interrupted,
			}

			assert.Equal(t, tc.expected, a.String())
		})
	}
}
