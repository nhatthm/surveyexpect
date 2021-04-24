# Expects for `AlecAivazis/survey`

[![GitHub Releases](https://img.shields.io/github/v/release/nhatthm/surveyexpect)](https://github.com/nhatthm/surveyexpect/releases/latest)
[![Build Status](https://github.com/nhatthm/surveyexpect/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/surveyexpect/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/surveyexpect/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/surveyexpect)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/surveyexpect)](https://goreportcard.com/report/github.com/nhatthm/surveyexpect)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/surveyexpect)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

**surveyexpect** is an Expect library for [AlecAivazis/survey/v2](https://github.com/AlecAivazis/survey)

## Prerequisites

- `Go >= 1.14`

## Install

```bash
go get github.com/nhatthm/surveyexpect
```

## Usage

### Supported Types

Type | Supported | Supported Actions
:--- | :---: | :---
`Confirm` | ✓ | <ul><li>Answer `yes`, `no` or a custom one</li><li>Interrupt (`^C`)</li><li>Ask for help</li></ul>
`Editor` | ✘ | __*There is no plan for support*__
`Input` | ✓ | <ul><li>Answer</li><li>No answer</li><li>Suggestions with navigation (Arrow Up `↑`, Arrow Down `↓`, Tab `⇆`, Esc `⎋`, Enter `⏎`)</li><li>Interrupt (`^C`)</li><li>Ask for help</li></ul>
`Multiline` | ✓ | <ul><li>Answer</li><li>No answer</li><li>Interrupt (`^C`)</li></ul>
`Multiselect` | ✓ | <ul><li>Type to filter, delete</li><li>Navigation (Move Up `↑`, Move Down `↓`, Select None `←`, Select All `→`, Tab `⇆`, Enter `⏎`)</li><li>Interrupt (`^C`)</li><li>Ask for help</li></ul>
`Password` | ✓ | <ul><li>Answer (+ check for `*`)</li><li>No answer</li><li>Interrupt (`^C`)</li><li>Ask for help</li></ul>
`Select` | ✓ | <ul><li>Type to filter, delete</li><li>Navigation (Move Up `↑`, Move Down `↓`, Tab `⇆`, Enter `⏎`)</li><li>Interrupt (`^C`)</li><li>Ask for help</li></ul>

### Expect

There are 2 steps:

Step 1: Create an expectation.

Call `surveyexpect.Expect()`

```go
s := surveyexpect.Expect(func(s *surveyexpect.Survey) {
    s.ExpectPassword("Enter a password:").
        Answer("secret")
})(t) // t is *testing.T
```

Step 2: Run it.

Important: Use the `stdio` arg and inject it into the `survey.Prompt` otherwise it won't work. 

```go
s.Start(func(stdio terminal.Stdio)) {
    // For example
    p := &survey.Password{Message: "Enter a password:"}
    var answer string
    err := survey.AskOne(p, &answer, surveyexpect.WithStdio(stdio))

    // Asserts.
    assert.Equal(t, "123456", answer)
    assert.NoError(t, err)
})
```

## Examples

```go
package mypackage_test

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/nhatthm/surveyexpect"
	"github.com/stretchr/testify/assert"
)

func TestMyPackage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		expectSurvey   surveyexpect.Expector
		expectedAnswer string
		expectedError  string
	}{
		{
			scenario: "empty answer",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter a password:").
					Answer("")
			}),
		},
		{
			scenario: "password without help",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter a password:").
					Answer("secret")
			}),
			expectedAnswer: "secret",
		},

		{
			scenario: "input is interrupted",
			expectSurvey: surveyexpect.Expect(func(s *surveyexpect.Survey) {
				s.ExpectPassword("Enter a password:").
					Interrupt()
			}),
			expectedError: "interrupt",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.scenario, func(t *testing.T) {
			t.Parallel()

			p := &survey.Password{Message: "Enter a password:"}

			// Start the survey.
			tc.expectSurvey(t).Start(func(stdio terminal.Stdio) {
				// Run your logic here.
				// For example.
				var answer string
				err := survey.AskOne(p, &answer, surveyexpect.WithStdio(stdio))

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
```

You can find more examples in the tests of this library:
- Confirm: https://github.com/nhatthm/surveyexpect/blob/master/confirm_test.go
- Input: https://github.com/nhatthm/surveyexpect/blob/master/input_test.go
- Multiline: https://github.com/nhatthm/surveyexpect/blob/master/multiline_test.go
- Multiselect: https://github.com/nhatthm/surveyexpect/blob/master/multiselect_test.go
- Password: https://github.com/nhatthm/surveyexpect/blob/master/password_test.go
- Select: https://github.com/nhatthm/surveyexpect/blob/master/select_test.go

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
