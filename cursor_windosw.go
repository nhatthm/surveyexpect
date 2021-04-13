// +build windows

package surveymock

func waitForCursor(c Console) error {
	<-time.After(ReactionTime)

	return nil
}
