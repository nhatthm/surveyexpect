// +build windows

package surveyexpect

func waitForCursor(c Console) error {
	<-time.After(ReactionTime)

	return nil
}
