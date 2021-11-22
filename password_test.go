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

func TestPasswordPrompt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		expectSurvey   surveyexpect.Expector
		message        string
		help           string
		showHelp       bool
		options        []survey.AskOpt
		expectedAnswer string
		expectedError  string
	}{
		{
			scenario: "no answer sends an empty line",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter an empty password:")
			}),
			message: "Enter an empty password:",
		},
		{
			scenario: "empty answer",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter an empty password:").
					Answer("")
			}),
			message: "Enter an empty password:",
		},
		{
			scenario: "password without help",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter a password:").
					Answer("secret")
			}),
			message:        "Enter a password:",
			expectedAnswer: "secret",
		},
		{
			scenario: "password with visible help and do not ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter a password: [? for help]").
					Answer("secret")
			}),
			message:        "Enter a password:",
			help:           "It is your secret",
			showHelp:       true,
			expectedAnswer: "secret",
		},
		{
			scenario: "password with visible help and ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter a password: [? for help]").
					ShowHelp("It is your secret")

				s.ExpectPassword("Enter a password:").
					Answer("secret")
			}),
			message:        "Enter a password:",
			help:           "It is your secret",
			showHelp:       true,
			expectedAnswer: "secret",
		},
		{
			scenario: "password with invisible help and do not ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter a password:").
					Answer("secret")
			}),
			message:        "Enter a password:",
			help:           "It is your secret",
			expectedAnswer: "secret",
		},
		{
			scenario: "password with invisible help and ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter a password:").
					ShowHelp("It is your secret")

				s.ExpectPassword("Enter a password:").
					Answer("secret")
			}),
			message:        "Enter a password:",
			help:           "It is your secret",
			expectedAnswer: "secret",
		},
		{
			scenario: "input is interrupted",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter a password:").
					Times(10). // Times will be discarded due to the interruption.
					Interrupt()
			}),
			message:       "Enter a password:",
			expectedError: "interrupt",
		},
		{
			scenario: "input is invalid",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter a password:").
					Answer("\033X").
					Interrupted()
			}),
			message:       "Enter a password:",
			expectedError: `unexpected escape sequence from terminal: ['\x1b' 'X']`,
		},
		{
			scenario: "answer is required",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				// Be asked for 5 times without giving up the password.
				s.ExpectPassword("Enter a password:").
					Times(5)

				// Finally, input the password.
				s.ExpectPassword("Enter a password:").
					Answer("secret")
			}),
			options: []survey.AskOpt{
				survey.WithValidator(survey.Required),
			},
			message:        "Enter a password:",
			expectedAnswer: "secret",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Prepare the survey.
			s := tc.expectSurvey(t)
			p := &survey.PasswordTemplateData{
				Password: survey.Password{Message: tc.message, Help: tc.help},
				ShowHelp: tc.showHelp,
			}

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

func TestPasswordPrompt_NoHelpButStillExpect(t *testing.T) {
	t.Parallel()

	testingT := T()
	s := surveyexpect.Expect(func(s *surveyexpect.Survey) {
		s.WithTimeout(50 * time.Millisecond)

		s.ExpectPassword("Enter a password:").
			ShowHelp("It is your secret")
	})(testingT)

	expectedAnswer := "?"
	expectedError := "there are remaining expectations that were not met:\n\nExpect : Password Prompt\nMessage: \"Enter a password:\"\nAnswer : ?\n"

	p := &survey.Password{Message: "Enter a password:"}

	// Start the survey.
	s.Start(func(stdio terminal.Stdio) {
		var answer string
		err := survey.AskOne(p, &answer, options.WithStdio(stdio))

		assert.Equal(t, expectedAnswer, answer)
		assert.NoError(t, err)
	})

	assert.EqualError(t, s.ExpectationsWereMet(), expectedError)

	t.Log(testingT.LogString())
}

func TestPasswordPrompt_SurveyInterrupted(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		expectSurvey  surveyexpect.Expector
		expectedError string
	}{
		{
			scenario: "interrupt",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter your username:").Interrupt()
			}),
			expectedError: "interrupt",
		},
		{
			scenario: "invalid input",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter your username:").
					Answer("\033X").
					Interrupted()
			}),
			expectedError: `unexpected escape sequence from terminal: ['\x1b' 'X']`,
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
					Prompt: &survey.Password{Message: "Enter your username:"},
				},
				{
					Name:   "password",
					Prompt: &survey.Password{Message: "Enter your password:"},
				},
			}

			expectedResult := map[string]interface{}{
				"username": "old username",
				"password": "old password",
			}

			// Start the survey.
			s.Start(func(stdio terminal.Stdio) {
				result := map[string]interface{}{
					"username": "old username",
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
