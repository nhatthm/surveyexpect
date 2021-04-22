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

func TestSelectPrompt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		expectSurvey   surveyexpect.Expector
		help           string
		showHelp       bool
		options        []string
		expectedAnswer string
		expectedError  string
	}{
		{
			scenario: "enter without taking any other actions",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectSelect("Select a country").
					Enter()
			}),
			expectedAnswer: "France",
		},
		{
			scenario: "with help and ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectSelect("Select a country  [Use arrows to move, type to filter, ? for more help]").
					ShowHelp("Your favorite country").
					Enter()
			}),
			help:           "Your favorite country",
			showHelp:       true,
			expectedAnswer: "France",
		},
		{
			scenario: "input is interrupted",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectSelect("Select a country").
					Interrupt()
			}),
			expectedError: "interrupt",
		},
		{
			scenario: "input is invalid",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectSelect("Select a country").
					Type("\033X")
			}),
			expectedError: `Unexpected Escape Sequence: ['\x1b' 'X']`,
		},
		{
			scenario: "navigation",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectSelect("Select a country").
					Type("United").
					ExpectOptions(
						"> United Kingdom",
						"United States",
					).
					Delete(6).
					Tab(2).MoveDown().MoveUp(4).
					ExpectOptions(
						"Germany",
						"Malaysia",
						"Singapore",
						"Thailand",
						"United Kingdom",
						"United States",
						"> Vietnam",
					).
					Enter()
			}),
			expectedAnswer: "Vietnam",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Prepare the survey.
			s := tc.expectSurvey(t)
			p := &survey.SelectTemplateData{
				Select: survey.Select{
					Message: "Select a country",
					Help:    tc.help,
					Options: []string{
						"France",
						"Germany",
						"Malaysia",
						"Singapore",
						"Thailand",
						"United Kingdom",
						"United States",
						"Vietnam",
					},
				},
				ShowHelp: tc.showHelp,
			}

			// Start the survey.
			s.Start(func(stdio terminal.Stdio) {
				var answer string
				err := survey.AskOne(p, &answer, options.WithStdio(stdio))

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

func TestSelectPrompt_NoHelpButStillExpect(t *testing.T) {
	t.Parallel()

	testingT := T()
	s := surveyexpect.Expect(func(s *surveyexpect.Survey) {
		s.WithTimeout(50 * time.Millisecond)

		s.ExpectSelect("Select a country").
			ShowHelp("Your favorite country")
	})(testingT)

	expectedError := "there are remaining expectations that were not met:\n\nExpect : Select Prompt\nMessage: \"Select a country\"\npress \"?\""

	p := &survey.Select{
		Message: "Select a country",
		Options: []string{
			"option 1",
			"option 2",
		},
	}

	// Start the survey.
	s.Start(func(stdio terminal.Stdio) {
		var answer string
		err := survey.AskOne(p, &answer, options.WithStdio(stdio))

		assert.Empty(t, answer)
		assert.NoError(t, err)
	})

	assert.EqualError(t, s.ExpectationsWereMet(), expectedError)

	t.Log(testingT.LogString())
}

func TestSelectPrompt_SurveyInterrupted(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		expectSurvey  surveyexpect.Expector
		expectedError string
	}{
		{
			scenario: "interrupt",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectSelect("Select a country").
					Interrupt()
			}),
			expectedError: "interrupt",
		},
		{
			scenario: "invalid input",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectSelect("Select a country").
					Type("\033X")
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
					Name:   "country",
					Prompt: &survey.Select{Message: "Select a country", Options: []string{"Germany", "Vietnam"}},
				},
				{
					Name:   "transport",
					Prompt: &survey.Select{Message: "Select a transport", Options: []string{"Train", "Bus"}},
				},
			}

			expectedResult := map[string]interface{}{
				"country":   "Vietnam",
				"transport": "Bus",
			}

			// Start the survey.
			s.Start(func(stdio terminal.Stdio) {
				result := map[string]interface{}{
					"country":   "Vietnam",
					"transport": "Bus",
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
