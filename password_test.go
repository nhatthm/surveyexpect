package surveymock_test

import (
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/surveymock"
)

func TestPassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		mockSurvey     surveymock.Mocker
		message        string
		help           string
		showHelp       bool
		options        []survey.AskOpt
		expectedAnswer string
		expectedError  string
	}{
		{
			scenario: "no answer sends an empty line",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectPassword("Enter an empty password:")
			}),
			message: "Enter an empty password:",
		},
		{
			scenario: "empty answer",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectPassword("Enter an empty password:").
					Answer("")
			}),
			message: "Enter an empty password:",
		},
		{
			scenario: "password without help",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectPassword("Enter a password:").
					Answer("secret")
			}),
			message:        "Enter a password:",
			expectedAnswer: "secret",
		},
		{
			scenario: "password with visible help and do not ask for it",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
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
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
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
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectPassword("Enter a password:").
					Answer("secret")
			}),
			message:        "Enter a password:",
			help:           "It is your secret",
			expectedAnswer: "secret",
		},
		{
			scenario: "password with invisible help and ask for it",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
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
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectPassword("Enter a password:").
					Times(10). // Times will be discarded due to the interruption.
					Interrupt()
			}),
			message:       "Enter a password:",
			expectedError: "interrupt",
		},
		{
			scenario: "input is invalid",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectPassword("Enter a password:").
					Answer("\033X").
					Interrupted()
			}),
			message:       "Enter a password:",
			expectedError: `Unexpected Escape Sequence: ['\x1b' 'X']`,
		},
		{
			scenario: "answer is required",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
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
			s := tc.mockSurvey(t)
			p := &survey.PasswordTemplateData{
				Password: survey.Password{Message: tc.message, Help: tc.help},
				ShowHelp: tc.showHelp,
			}

			// Start the survey.
			s.Start(func(stdio terminal.Stdio) {
				tc.options = append(tc.options, surveymock.WithStdio(stdio))

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

func TestPassword_NoHelpButStillExpect(t *testing.T) {
	t.Parallel()

	testingT := T()
	s := surveymock.Mock(func(s *surveymock.Survey) {
		s.WithTimeout(10 * time.Millisecond)

		s.ExpectPassword("Enter a password:").
			ShowHelp("It is your secret")
	})(testingT)

	expectedAnswer := "?"
	expectedError := "there are remaining expectations that were not met:\n\nType   : Password\nMessage: \"Enter a password:\"\nAnswer : ?\n"

	p := &survey.Password{Message: "Enter a password:"}

	// Start the survey.
	s.Start(func(stdio terminal.Stdio) {
		var answer string
		err := survey.AskOne(p, &answer, surveymock.WithStdio(stdio))

		assert.Equal(t, expectedAnswer, answer)
		assert.NoError(t, err)
	})

	assert.EqualError(t, s.ExpectationsWereMet(), expectedError)

	t.Log(testingT.LogString())
}

func TestPassword_SurveyInterrupted(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		mockSurvey    surveymock.Mocker
		expectedError string
	}{
		{
			scenario: "interrupt",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectPassword("Enter your username:").Interrupt()
			}),
			expectedError: "interrupt",
		},
		{
			scenario: "invalid input",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectPassword("Enter your username:").
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
			s := tc.mockSurvey(testingT)

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

			expectedResult := map[string]string{
				"username": "old username",
				"password": "old password",
			}

			// Start the survey.
			s.Start(func(stdio terminal.Stdio) {
				result := map[string]string{
					"username": "old username",
					"password": "old password",
				}
				err := survey.Ask(questions, &result, surveymock.WithStdio(stdio))

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
