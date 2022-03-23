package surveyexpect

import (
	"testing"

	"github.com/Netflix/go-expect"
	pseudotty "github.com/creack/pty"
	"github.com/hinshun/vt10x"
	"github.com/stretchr/testify/require"
)

func TestWaitForCursor(t *testing.T) {
	t.Parallel()

	pty, tty, err := pseudotty.Open()
	require.NoError(t, err)

	term := vt10x.New(vt10x.WithWriter(tty))

	console, err := expect.NewConsole(expect.WithStdin(pty), expect.WithStdout(term), expect.WithCloser(pty, tty))
	require.NoError(t, err)

	_ = console.Close() // nolint: errcheck
	_ = tty.Close()     // nolint: errcheck

	err = waitForCursor(console)
	require.Error(t, err)
}

func TestWaitForCursorTwice(t *testing.T) {
	t.Parallel()

	pty, tty, err := pseudotty.Open()
	require.NoError(t, err)

	term := vt10x.New(vt10x.WithWriter(tty))

	console, err := expect.NewConsole(expect.WithStdin(pty), expect.WithStdout(term), expect.WithCloser(pty, tty))
	require.NoError(t, err)

	_ = console.Close() // nolint: errcheck
	_ = tty.Close()     // nolint: errcheck

	err = waitForCursorTwice(console)
	require.Error(t, err)
}
