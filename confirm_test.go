package surveyexpect_test

import (
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/surveyexpect"
	"github.com/nhatthm/surveyexpect/options"
)

func TestConfirm(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		expectSurvey   surveyexpect.Expector
		defaultValue   bool
		help           string
		showHelp       bool
		options        []survey.AskOpt
		expectedAnswer bool
		expectedError  string
	}{
		{
			scenario: "no answer sends an empty answer (default: false)",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt?")
			}),
		},
		{
			scenario: "empty answer (default: false)",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt?").
					Answer("")
			}),
		},
		{
			scenario: "no answer sends an empty answer (default: true)",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt?")
			}),
			defaultValue:   true,
			expectedAnswer: true,
		},
		{
			scenario: "empty answer (default: true)",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt?").
					Answer("")
			}),
			defaultValue:   true,
			expectedAnswer: true,
		},
		{
			scenario: "confirm without help (yes)",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt?").Yes()
			}),
			defaultValue:   false,
			expectedAnswer: true,
		},
		{
			scenario: "confirm without help (no)",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt?").No()
			}),
			defaultValue:   true,
			expectedAnswer: false,
		},
		{
			scenario: "confirm with visible help and do not ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt? [? for help]").Yes()
			}),
			help:           "This is a helpful help",
			showHelp:       true,
			expectedAnswer: true,
		},
		{
			scenario: "confirm with visible help and ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt? [? for help]").
					ShowHelp("This is a helpful help")

				s.ExpectConfirm("ConfirmPrompt?").Yes()
			}),
			help:           "This is a helpful help",
			showHelp:       true,
			expectedAnswer: true,
		},
		{
			scenario: "confirm with invisible help and do not ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt?").Yes()
			}),
			help:           "This is a helpful help",
			expectedAnswer: true,
		},
		{
			scenario: "confirm with invisible help and ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt?").
					ShowHelp("This is a helpful help")

				s.ExpectConfirm("ConfirmPrompt?").Yes()
			}),
			help:           "This is a helpful help",
			expectedAnswer: true,
		},
		{
			scenario: "input is interrupted",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt?").
					Interrupt()
			}),
			expectedError: "interrupt",
		},
		{
			scenario: "input contains invalid character",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt?").
					Answer("\033X").
					Interrupted()
			}),
			expectedError: `Unexpected Escape Sequence: ['\x1b' 'X']`,
		},
		{
			scenario: "input is invalid",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt?").
					Answer("not a yes or no")

				s.ExpectConfirm("ConfirmPrompt?").Yes()
			}),
			expectedAnswer: true,
		},
		{
			scenario: "required does not affect confirm",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("ConfirmPrompt?")
			}),
			options: []survey.AskOpt{
				survey.WithValidator(survey.Required),
			},
			expectedAnswer: false,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Prepare the survey.
			s := tc.expectSurvey(t)
			p := &survey.ConfirmTemplateData{
				Confirm:  survey.Confirm{Message: "ConfirmPrompt?", Help: tc.help, Default: tc.defaultValue},
				ShowHelp: tc.showHelp,
			}

			// Start the survey.
			s.Start(func(stdio terminal.Stdio) {
				tc.options = append(tc.options, options.WithStdio(stdio))

				var answer bool
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

func TestConfirm_NoHelpButStillExpect(t *testing.T) {
	t.Parallel()

	testingT := T()
	s := surveyexpect.Expect(func(s *surveyexpect.Survey) {
		s.WithTimeout(50 * time.Millisecond)

		s.ExpectConfirm("ConfirmPrompt?").
			ShowHelp("It is your secret")
	})(testingT)

	expectedError := "there are remaining expectations that were not met:\n\nExpect : Confirm Prompt\nMessage: \"ConfirmPrompt?\"\nAnswer : ?\n"

	p := &survey.Confirm{Message: "ConfirmPrompt?"}

	// Start the survey.
	s.Start(func(stdio terminal.Stdio) {
		var answer bool
		err := survey.AskOne(p, &answer, options.WithStdio(stdio))

		assert.False(t, answer)
		assert.NoError(t, err)
	})

	assert.EqualError(t, s.ExpectationsWereMet(), expectedError)

	t.Log(testingT.LogString())
}

func TestConfirm_SurveyInterrupted(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		expectSurvey  surveyexpect.Expector
		expectedError string
	}{
		{
			scenario: "interrupt",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to save your password?").Interrupt()
			}),
			expectedError: "interrupt",
		},
		{
			scenario: "invalid input",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectConfirm("Do you want to save your password?").
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
					Name:   "username",
					Prompt: &survey.Confirm{Message: "Do you want to save your password?"},
				},
				{
					Name:   "password",
					Prompt: &survey.Password{Message: "Enter your password:"},
				},
			}

			expectedResult := map[string]interface{}{
				"confirm":  true,
				"password": "old password",
			}

			// Start the survey.
			s.Start(func(stdio terminal.Stdio) {
				result := map[string]interface{}{
					"confirm":  true,
					"password": "old password",
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
