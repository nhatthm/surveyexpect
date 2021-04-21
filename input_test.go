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

func TestInputPrompt(t *testing.T) {
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
				s.ExpectInput("Enter an empty username:")
			}),
			message: "Enter an empty username:",
		},
		{
			scenario: "empty answer",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectInput("Enter an empty username:").
					Answer("")
			}),
			message: "Enter an empty username:",
		},
		{
			scenario: "username without help",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectInput("Enter a username:").
					Answer("secret")
			}),
			message:        "Enter a username:",
			expectedAnswer: "secret",
		},
		{
			scenario: "username with visible help and do not ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectInput("Enter a username: [? for help]").
					Answer("secret")
			}),
			message:        "Enter a username:",
			help:           "It is your email",
			showHelp:       true,
			expectedAnswer: "secret",
		},
		{
			scenario: "username with visible help and ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectInput("Enter a username: [? for help]").
					ShowHelp("It is your email")

				s.ExpectInput("Enter a username:").
					Answer("secret")
			}),
			message:        "Enter a username:",
			help:           "It is your email",
			showHelp:       true,
			expectedAnswer: "secret",
		},
		{
			scenario: "username with invisible help and do not ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectInput("Enter a username:").
					Answer("secret")
			}),
			message:        "Enter a username:",
			help:           "It is your email",
			expectedAnswer: "secret",
		},
		{
			scenario: "username with invisible help and ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectInput("Enter a username:").
					ShowHelp("It is your email")

				s.ExpectInput("Enter a username:").
					Answer("secret")
			}),
			message:        "Enter a username:",
			help:           "It is your email",
			expectedAnswer: "secret",
		},
		{
			scenario: "input is interrupted",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectInput("Enter a username:").
					Times(10). // Times will be discarded due to the interruption.
					Interrupt()
			}),
			message:       "Enter a username:",
			expectedError: "interrupt",
		},
		{
			scenario: "input is invalid",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectInput("Enter a username:").
					Answer("\033X").
					Interrupted()
			}),
			message:       "Enter a username:",
			expectedError: `Unexpected Escape Sequence: ['\x1b' 'X']`,
		},
		{
			scenario: "answer is required",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				// Be asked for 5 times without giving up the username.
				s.ExpectInput("Enter a username:").
					Times(5)

				// Finally, input the username.
				s.ExpectInput("Enter a username:").
					Answer("secret")
			}),
			options: []survey.AskOpt{
				survey.WithValidator(survey.Required),
			},
			message:        "Enter a username:",
			expectedAnswer: "secret",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Prepare the survey.
			s := tc.expectSurvey(t)
			p := &survey.InputTemplateData{
				Input:    survey.Input{Message: tc.message, Help: tc.help},
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

func TestInputPrompt_NoHelpButStillExpect(t *testing.T) {
	t.Parallel()

	testingT := T()
	s := surveyexpect.Expect(func(s *surveyexpect.Survey) {
		s.WithTimeout(50 * time.Millisecond)

		s.ExpectInput("Enter a username:").
			ShowHelp("It is your email")
	})(testingT)

	expectedAnswer := "?"
	expectedError := "there are remaining expectations that were not met:\n\nExpect : Input Prompt\nMessage: \"Enter a username:\"\nAnswer : ?\n"

	p := &survey.Input{Message: "Enter a username:"}

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

func TestInputPrompt_SurveyInterrupted(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		expectSurvey  surveyexpect.Expector
		expectedError string
	}{
		{
			scenario: "interrupt",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectInput("Enter your username:").Interrupt()
			}),
			expectedError: "interrupt",
		},
		{
			scenario: "invalid input",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectInput("Enter your username:").
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
					Prompt: &survey.Input{Message: "Enter your username:"},
				},
				{
					Name:   "password",
					Prompt: &survey.Input{Message: "Enter your password:"},
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

func TestInputPrompt_AskForSuggestions(t *testing.T) {
	t.Parallel()

	s := surveyexpect.Expect(func(s *surveyexpect.Survey) {
		s.ExpectInput("Enter username:").
			Type("joh").Tab().
			ExpectSuggestions(
				"> john.doe",
				"john.lennon",
				"john.legend",
				"john.mayor",
				"john.micheal",
				"john.nguyen",
				"john.pierre",
			).
			Tab().Tab().MoveUp().MoveUp().MoveDown().
			ExpectSuggestions(
				"john.doe",
				"> john.lennon",
				"john.legend",
				"john.mayor",
				"john.micheal",
				"john.nguyen",
				"john.pierre",
			).
			Esc().Tab().Type("n").Tab().MoveUp().
			ExpectSuggestions(
				"john.lennon",
				"john.legend",
				"john.mayor",
				"john.micheal",
				"john.nguyen",
				"john.pierre",
				"> johnny",
			).
			Enter()
	})(t)

	p := &survey.Input{
		Message: "Enter username:",
		Suggest: func(string) []string {
			return []string{
				"john.doe",
				"john.lennon",
				"john.legend",
				"john.mayor",
				"john.micheal",
				"john.nguyen",
				"john.pierre",
				"johnny",
			}
		},
	}

	expectedAnswer := `johnny`

	// Start the survey.
	s.Start(func(stdio terminal.Stdio) {
		var answer string
		err := survey.AskOne(p, &answer, options.WithStdio(stdio))

		assert.Equal(t, expectedAnswer, answer)
		assert.NoError(t, err)
	})
}

func TestInputPrompt_AskForSuggestionsButThereIsNone(t *testing.T) {
	t.Parallel()

	testingT := T()
	s := surveyexpect.Expect(func(s *surveyexpect.Survey) {
		s.WithTimeout(100 * time.Millisecond)
		s.ExpectInput("Enter username:").
			Type("john").Tab().
			ExpectSuggestions(
				"> john.doe",
				"john.lennon",
				"john.legend",
				"john.mayor",
				"john.micheal",
				"john.nguyen",
				"john.pierre",
			).
			Enter()
	})(testingT)

	p := &survey.Input{Message: "Enter username:"}

	// Start the survey.
	s.Start(func(stdio terminal.Stdio) {
		var answer string
		err := survey.AskOne(p, &answer, options.WithStdio(stdio))

		assert.Empty(t, answer)
		assert.NoError(t, err)
	})
}
