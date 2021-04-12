package options

import (
	"bytes"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type buffer struct {
	bytes.Buffer
}

func (b *buffer) Fd() uintptr {
	return 0
}

func TestWithStdio(t *testing.T) {
	t.Parallel()

	buf := &buffer{}
	stdio := terminal.Stdio{
		In:  buf,
		Out: buf,
		Err: buf,
	}

	result := &survey.AskOptions{}
	err := WithStdio(stdio)(result)
	require.NoError(t, err)

	assert.Equal(t, buf, result.Stdio.In)
	assert.Equal(t, buf, result.Stdio.Out)
	assert.Equal(t, buf, result.Stdio.Err)
}
