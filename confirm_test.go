package surveymock_test

import (
	"testing"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/surveymock"
	"github.com/nhatthm/surveymock/options"
)

func TestConfirm(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		mockSurvey     surveymock.Mocker
		defaultValue   bool
		help           string
		showHelp       bool
		options        []survey.AskOpt
		expectedAnswer bool
		expectedError  string
	}{
		{
			scenario: "no answer sends an empty answer (default: false)",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm?")
			}),
		},
		{
			scenario: "empty answer (default: false)",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm?").
					Answer("")
			}),
		},
		{
			scenario: "no answer sends an empty answer (default: true)",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm?")
			}),
			defaultValue:   true,
			expectedAnswer: true,
		},
		{
			scenario: "empty answer (default: true)",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm?").
					Answer("")
			}),
			defaultValue:   true,
			expectedAnswer: true,
		},
		{
			scenario: "confirm without help (yes)",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm?").Yes()
			}),
			defaultValue:   false,
			expectedAnswer: true,
		},
		{
			scenario: "confirm without help (no)",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm?").No()
			}),
			defaultValue:   true,
			expectedAnswer: false,
		},
		{
			scenario: "confirm with visible help and do not ask for it",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm? [? for help]").Yes()
			}),
			help:           "This is a helpful help",
			showHelp:       true,
			expectedAnswer: true,
		},
		{
			scenario: "confirm with visible help and ask for it",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm? [? for help]").
					ShowHelp("This is a helpful help")

				s.ExpectConfirm("Confirm?").Yes()
			}),
			help:           "This is a helpful help",
			showHelp:       true,
			expectedAnswer: true,
		},
		{
			scenario: "confirm with invisible help and do not ask for it",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm?").Yes()
			}),
			help:           "This is a helpful help",
			expectedAnswer: true,
		},
		{
			scenario: "confirm with invisible help and ask for it",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm?").
					ShowHelp("This is a helpful help")

				s.ExpectConfirm("Confirm?").Yes()
			}),
			help:           "This is a helpful help",
			expectedAnswer: true,
		},
		{
			scenario: "input is interrupted",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm?").
					Interrupt()
			}),
			expectedError: "interrupt",
		},
		{
			scenario: "input contains invalid character",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm?").
					Answer("\033X").
					Interrupted()
			}),
			expectedError: `Unexpected Escape Sequence: ['\x1b' 'X']`,
		},
		{
			scenario: "input is invalid",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm?").
					Answer("not a yes or no")

				s.ExpectConfirm("Confirm?").Yes()
			}),
			expectedAnswer: true,
		},
		{
			scenario: "required does not affect confirm",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Confirm?")
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
			s := tc.mockSurvey(t)
			p := &survey.ConfirmTemplateData{
				Confirm:  survey.Confirm{Message: "Confirm?", Help: tc.help, Default: tc.defaultValue},
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
	s := surveymock.Mock(func(s *surveymock.Survey) {
		s.WithTimeout(10 * time.Millisecond)

		s.ExpectConfirm("Confirm?").
			ShowHelp("It is your secret")
	})(testingT)

	expectedError := "there are remaining expectations that were not met:\n\nType   : Confirm\nMessage: \"Confirm?\"\nAnswer : ?\n"

	p := &survey.Confirm{Message: "Confirm?"}

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
		mockSurvey    surveymock.Mocker
		expectedError string
	}{
		{
			scenario: "interrupt",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectConfirm("Do you want to save your password?").Interrupt()
			}),
			expectedError: "interrupt",
		},
		{
			scenario: "invalid input",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
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
			s := tc.mockSurvey(testingT)

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
