package cobra

import (
	"io"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"

	"github.com/nhatthm/surveymock/options"
)

// StdioProvider is a wrapper around *cobra.Command to provide stdin, stdout and stderr to survey.
type StdioProvider interface {
	OutOrStdout() io.Writer
	ErrOrStderr() io.Writer
	InOrStdin() io.Reader
}

// WithStdioProvider configures stdio for prompt.
func WithStdioProvider(p StdioProvider) survey.AskOpt {
	in, ok := p.InOrStdin().(terminal.FileReader)
	if !ok {
		return configureNothing
	}

	out, ok := p.OutOrStdout().(terminal.FileWriter)
	if !ok {
		return configureNothing
	}

	return options.WithStdio(terminal.Stdio{
		In:  in,
		Out: out,
		Err: p.ErrOrStderr(),
	})
}

func configureNothing(*survey.AskOptions) error {
	return nil
}
