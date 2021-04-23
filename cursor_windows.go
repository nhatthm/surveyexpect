// +build windows

package surveyexpect

func waitForCursor(c Console) error {
	<-WaitForReaction()

	return nil
}
