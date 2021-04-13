// +build darwin

package surveymock

func waitForCursor(c Console) error {
	// ANSI escape sequence for DSR - Device Status Report
	// https://en.wikipedia.org/wiki/ANSI_escape_code#CSI_sequences
	_, err := c.ExpectString("\x1b[6n")

	return err
}
