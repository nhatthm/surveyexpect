package surveyexpect

import (
	"testing"

	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/require"
)

func TestWaitForCursor(t *testing.T) {
	t.Parallel()

	console, _, err := vt10x.NewVT10XConsole()
	require.NoError(t, err)

	_ = console.Tty().Close() // nolint: errcheck
	_ = console.Close()       // nolint: errcheck

	err = waitForCursor(console)
	require.Error(t, err)
}

func TestWaitForCursorTwice(t *testing.T) {
	t.Parallel()

	console, _, err := vt10x.NewVT10XConsole()
	require.NoError(t, err)

	_ = console.Tty().Close() // nolint: errcheck
	_ = console.Close()       // nolint: errcheck

	err = waitForCursorTwice(console)
	require.Error(t, err)
}
