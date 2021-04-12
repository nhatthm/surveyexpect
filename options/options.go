package options

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

// WithStdio sets survey.WithStdio for a prompt.
func WithStdio(stdio terminal.Stdio) survey.AskOpt {
	return survey.WithStdio(stdio.In, stdio.Out, stdio.Err)
}
