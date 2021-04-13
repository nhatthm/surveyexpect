package surveymock

func waitForCursorTwice(c Console) error {
	if err := waitForCursor(c); err != nil {
		return err
	}

	return waitForCursor(c)
}
