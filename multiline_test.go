package surveyexpect_test

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/surveyexpect"
	"github.com/nhatthm/surveyexpect/options"
)

func TestMultilinePrompt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		expectSurvey   surveyexpect.Expector
		help           string
		showHelp       bool
		options        []survey.AskOpt
		expectedAnswer string
		expectedError  string
	}{
		{
			scenario: "no answer sends an empty line",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectMultiline("Enter your comment")
			}),
		},
		{
			scenario: "empty answer",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectMultiline("Enter your comment").
					Answer("")
			}),
		},
		{
			scenario: "input is interrupted",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectMultiline("Enter your comment").
					Times(10). // Times will be discarded due to the interruption.
					Interrupt()
			}),
			expectedError: "interrupt",
		},
		{
			scenario: "input is invalid",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectMultiline("Enter your comment").
					Answer("\033X").
					Interrupted()
			}),
			expectedError: `Unexpected Escape Sequence: ['\x1b' 'X']`,
		},
		{
			scenario: "answer is required",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				// Be asked for 5 times without giving up the password.
				s.ExpectMultiline("Enter your comment").
					Times(5)

				// Finally, input the password.
				s.ExpectMultiline("Enter your comment").
					Answer("this is a multiline\ncomment\n\nend")
			}),
			options: []survey.AskOpt{
				survey.WithValidator(survey.Required),
			},
			expectedAnswer: "this is a multiline\ncomment\n\nend",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Prepare the survey.
			s := tc.expectSurvey(t)
			p := &survey.Multiline{Message: "Enter your comment"}

			// Start the survey.
			s.Start(func(stdio terminal.Stdio) {
				tc.options = append(tc.options, options.WithStdio(stdio))

				var answer string
				err := survey.AskOne(p, &answer, tc.options...)

				assert.Equal(t, tc.expectedAnswer, answer)

				if tc.expectedError == "" {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.expectedError)
				}
			})
		})
	}
}

func TestMultilinePrompt_SurveyInterrupted(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		expectSurvey  surveyexpect.Expector
		expectedError string
	}{
		{
			scenario: "interrupt",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectMultiline("Enter your message:").Interrupt()
			}),
			expectedError: "interrupt",
		},
		{
			scenario: "invalid input",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectMultiline("Enter your message:").
					Answer("\033X").
					Interrupted()
			}),
			expectedError: `Unexpected Escape Sequence: ['\x1b' 'X']`,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			testingT := T()
			s := tc.expectSurvey(testingT)

			questions := []*survey.Question{
				{
					Name:   "message",
					Prompt: &survey.Multiline{Message: "Enter your message:"},
				},
				{
					Name:   "status",
					Prompt: &survey.Multiline{Message: "Enter your status:"},
				},
			}

			expectedResult := map[string]interface{}{
				"message": "old message",
				"status":  "old status",
			}

			// Start the survey.
			s.Start(func(stdio terminal.Stdio) {
				result := map[string]interface{}{
					"message": "old message",
					"status":  "old status",
				}
				err := survey.Ask(questions, &result, options.WithStdio(stdio))

				assert.Equal(t, expectedResult, result)

				if tc.expectedError == "" {
					assert.NoError(t, err)
				} else {
					assert.EqualError(t, err, tc.expectedError)
				}
			})

			assert.Nil(t, s.ExpectationsWereMet())

			t.Log(testingT.LogString())
		})
	}
}
