package surveyexpect

// StringsExpect expects strings from console.
type StringsExpect []string

// Do runs the step.
func (e *StringsExpect) Do(c Console) error {
	for _, s := range *e {
		if _, err := c.ExpectString(s); err != nil {
			return err
		}
	}

	return nil
}

// String represents the answer as a string.
func (e *StringsExpect) String() string {
	return "TODO"
}

func expectStrings(strings ...string) *StringsExpect {
	e := StringsExpect(strings)

	return &e
}
