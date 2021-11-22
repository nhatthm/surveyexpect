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

func TestMultiSelectPrompt(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		expectSurvey   surveyexpect.Expector
		help           string
		showHelp       bool
		options        []string
		expectedAnswer []string
		expectedError  string
	}{
		{
			scenario: "enter without taking any other actions",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectMultiSelect("Select destinations").
					Enter()
			}),
		},
		{
			scenario: "with help and ask for it",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectMultiSelect("Select destinations  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter, ? for more help]").
					ShowHelp("Your favorite countries").
					Select().
					Enter()
			}),
			help:           "Your favorite countries",
			showHelp:       true,
			expectedAnswer: []string{"France"},
		},
		{
			scenario: "input is interrupted",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectMultiSelect("Select destinations").
					ExpectOptions(
						"> [ ]  France",
						"[ ]  Germany",
						"[ ]  Malaysia",
						"[ ]  Singapore",
						"[ ]  Thailand",
						"[ ]  United Kingdom",
						"[ ]  United States",
					).
					Interrupt()
			}),
			expectedError: "interrupt",
		},
		{
			scenario: "input is invalid",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectMultiSelect("Select destinations").
					Type("\033X")
			}),
			expectedError: `unexpected escape sequence from terminal: ['\x1b' 'X']`,
		},
		{
			scenario: "navigation",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectMultiSelect("Select destinations").
					Type("United").Delete(2).
					ExpectOptions(
						"> [ ]  United Kingdom",
						"[ ]  United States",
					).
					SelectAll().
					Tab(2).MoveDown().MoveUp(4).
					ExpectOptions(
						"[ ]  Germany",
						"[ ]  Malaysia",
						"[ ]  Singapore",
						"[ ]  Thailand",
						"[x]  United Kingdom",
						"[x]  United States",
						"> [ ]  Vietnam",
					).
					MoveDown().Select().
					ExpectOptions(
						"[x]  France",
						"[ ]  Germany",
						"[ ]  Malaysia",
						"[ ]  Singapore",
						"[ ]  Thailand",
						"[x]  United Kingdom",
						"[x]  United States",
					).
					SelectNone().
					Select().Select().
					MoveUp().
					ExpectOptions(
						"[ ]  Germany",
						"[ ]  Malaysia",
						"[ ]  Singapore",
						"[ ]  Thailand",
						"[ ]  United Kingdom",
						"[ ]  United States",
						"> [ ]  Vietnam",
					).
					Select().
					Enter()
			}),
			expectedAnswer: []string{"Vietnam"},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			// Prepare the survey.
			s := tc.expectSurvey(t)
			p := &survey.MultiSelectTemplateData{
				MultiSelect: survey.MultiSelect{
					Message: "Select destinations",
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
				var answer []string
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

func TestMultiSelectPrompt_NoHelpButStillExpect(t *testing.T) {
	t.Parallel()

	testingT := T()
	s := surveyexpect.Expect(func(s *surveyexpect.Survey) {
		s.WithTimeout(50 * time.Millisecond)

		s.ExpectMultiSelect("Select destinations").
			ShowHelp("Your favorite countries").
			ExpectOptions(
				"> [ ]  option 1",
				"[ ]  option 2",
			)
	})(testingT)

	expectedError := `there are remaining expectations that were not met:

Expect : MultiSelect Prompt
Message: "Select destinations"
press "?" and see "Your favorite countries"
Expect a multiselect list:
> [ ]  option 1
  [ ]  option 2`

	p := &survey.MultiSelect{
		Message: "Select destinations",
		Options: []string{
			"option 1",
			"option 2",
		},
	}

	// Start the survey.
	s.Start(func(stdio terminal.Stdio) {
		var answer []string
		err := survey.AskOne(p, &answer, options.WithStdio(stdio))

		assert.Empty(t, answer)
		assert.NoError(t, err)
	})

	assert.EqualError(t, s.ExpectationsWereMet(), expectedError)

	t.Log(testingT.LogString())
}

func TestMultiSelectPrompt_SurveyInterrupted(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario      string
		expectSurvey  surveyexpect.Expector
		expectedError string
	}{
		{
			scenario: "interrupt",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectMultiSelect("Select destinations").
					ExpectOptions(
						"> [ ]  Germany",
						"[ ]  Vietnam",
					).
					Interrupt()
			}),
			expectedError: "interrupt",
		},
		{
			scenario: "invalid input",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectMultiSelect("Select destinations").
					Type("\033X")
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
					Name:   "countries",
					Prompt: &survey.MultiSelect{Message: "Select destinations", Options: []string{"Germany", "Vietnam"}},
				},
				{
					Name:   "transports",
					Prompt: &survey.MultiSelect{Message: "Select transports", Options: []string{"Train", "Bus"}},
				},
			}

			expectedResult := map[string]interface{}{
				"countries":  []string{"Vietnam"},
				"transports": []string{"Bus"},
			}

			// Start the survey.
			s.Start(func(stdio terminal.Stdio) {
				result := map[string]interface{}{
					"countries":  []string{"Vietnam"},
					"transports": []string{"Bus"},
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
