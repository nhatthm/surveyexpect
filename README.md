# Mock for `AlecAivazis/survey`

[![Build Status](https://github.com/nhatthm/surveymock/actions/workflows/test.yaml/badge.svg)](https://github.com/nhatthm/surveymock/actions/workflows/test.yaml)
[![codecov](https://codecov.io/gh/nhatthm/surveymock/branch/master/graph/badge.svg?token=eTdAgDE2vR)](https://codecov.io/gh/nhatthm/surveymock)
[![Go Report Card](https://goreportcard.com/badge/github.com/nhatthm/surveymock)](https://goreportcard.com/report/github.com/nhatthm/surveymock)
[![GoDevDoc](https://img.shields.io/badge/dev-doc-00ADD8?logo=go)](https://pkg.go.dev/github.com/nhatthm/surveymock)
[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

**surveymock** is a mock library for [AlecAivazis/survey/v2](https://github.com/AlecAivazis/survey)

## Prerequisites

- `Go >= 1.14`

## Install

```bash
go get github.com/nhatthm/surveymock
```

## Usage

### Supported Types

For now, it only supports `Confirm` and `Password`


### Mock

There are 2 steps:

Step 1: Create an expectation for the survey.

Call `surveymock.Mock()`

```go
survey := surveymock.Mock(func(s *surveymock.Survey) {
    s.ExpectPassword("Enter a password:").
        Answer("secret")
})(t) // t is *testing.T
```

Step 2: Run it.

Important: Use the `stdio` arg and inject it into the `survey.Prompt` otherwise it won't work. 

```go
survey.Start(func(stdio terminal.Stdio)) {
    // For example
    p := &survey.Password{Message: "Enter a password:"}
    var answer string
    err := survey.AskOne(p, &answer, surveymock.WithStdio(stdio))

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
	"github.com/nhatthm/surveymock"
	"github.com/stretchr/testify/assert"
)

func TestMyPackage(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		scenario       string
		mockSurvey     surveymock.Mocker
		expectedAnswer string
		expectedError  string
	}{
		{
			scenario: "empty answer",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectPassword("Enter a password:").
					Answer("")
			}),
		},
		{
			scenario: "password without help",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
				s.ExpectPassword("Enter a password:").
					Answer("secret")
			}),
			expectedAnswer: "secret",
		},

		{
			scenario: "input is interrupted",
			mockSurvey: surveymock.Mock(func(s *surveymock.Survey) {
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
			tc.mockSurvey(t).Start(func(stdio terminal.Stdio) {
				// Run your logic here.
				// For example.
				var answer string
				err := survey.AskOne(p, &answer, surveymock.WithStdio(stdio))

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

## Donation

If this project help you reduce time to develop, you can give me a cup of coffee :)

### Paypal donation

[![paypal](https://www.paypalobjects.com/en_US/i/btn/btn_donateCC_LG.gif)](https://www.paypal.com/donate/?hosted_button_id=PJZSGJN57TDJY)

&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;or scan this

<img src="https://user-images.githubusercontent.com/1154587/113494222-ad8cb200-94e6-11eb-9ef3-eb883ada222a.png" width="147px" />
