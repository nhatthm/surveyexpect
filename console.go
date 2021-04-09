package surveymock

import (
	"os"

	"github.com/Netflix/go-expect"
)

// Console is an interface wrapper around *expect.Console.
type Console interface {
	// terminal device.
	Tty() *os.File

	// pty.
	Fd() uintptr

	// Close closes Console's tty. Calling Close will unblock Expect and ExpectEOF.
	Close() error

	// Send writes string s to Console's tty.
	Send(s string) (int, error)

	// SendLine writes string s to Console's tty with a trailing newline.
	SendLine(s string) (int, error)

	// Expectf reads from the Console's tty until the provided formatted string
	// is read or an error occurs, and returns the buffer read by Console.
	Expectf(format string, args ...interface{}) (string, error)

	// ExpectString reads from Console's tty until the provided string is read or
	// an error occurs, and returns the buffer read by Console.
	ExpectString(s string) (string, error)

	// ExpectEOF reads from Console's tty until EOF or an error occurs, and returns
	// the buffer read by Console.  We also treat the PTSClosed error as an EOF.
	ExpectEOF() (string, error)

	// Expect reads from Console's tty until a condition specified from opts is
	// encountered or an error occurs, and returns the buffer read by console.
	// No extra bytes are read once a condition is met, so if a program isn't
	// expecting input yet, it will be blocked. Sends are queued up in tty's
	// internal buffer so that the next Expect will read the remaining bytes (i.e.
	// rest of prompt) as well as its conditions.
	Expect(opts ...expect.ExpectOpt) (string, error)
}
