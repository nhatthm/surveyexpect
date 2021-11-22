//go:build !windows
// +build !windows

package surveyexpect

func waitForCursor(c Console) error {
	// ANSI escape sequence for DSR - Device Status Report
	// https://en.wikipedia.org/wiki/ANSI_escape_code#CSI_sequences
	_, err := c.ExpectString("\x1b[6n")

	// Simulate human delay, this is to fix the issue on Linux when the addStep answers too fast and the response is not
	// caught on AlecAivazis/survey side.
	//
	// This is spotted by printing out a bunch of logs with micro time to see how everything reacts on both sides.
	// After rendering the question, the prompt asks for the cursor's size and location (ESC[6n) and expects to receive
	// `ESC[n;mR` in return before reading the answer. If the addStep answers too fast (so the answer will be in between
	// `ESC[n;mR` and reading answer), the prompt won't see the answer and hangs indefinitely.
	<-WaitForReaction()

	return err
}
