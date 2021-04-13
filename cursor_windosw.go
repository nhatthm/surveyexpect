// +build windows

package surveymock

func waitForCursor(c Console) error {
	<-time.After(2 * time.Millisecond)

	return nil
}
