package surveyexpect

import (
	"fmt"
	"regexp"
)

var (
	selectIndicatorRegex      = regexp.MustCompile(`^([^ ]\s+)(.*)`)
	multiselectIndicatorRegex = regexp.MustCompile(`^([^ ]\s+)(\[[x ]].*)`)
)

// SelectExpect expects a select list from console.
type SelectExpect []string

// Do runs the step.
func (e *SelectExpect) Do(c Console) error {
	for _, o := range *e {
		if _, err := c.ExpectString(o); err != nil {
			return err
		}
	}

	return nil
}

// String represents the answer as a string.
func (e *SelectExpect) String() string {
	var sb stringsBuilder

	sb.WriteLinef("Expect a select list:")
	writeSelectList(&sb, *e, selectIndicatorRegex)

	return sb.String()
}

func expectSelect(options ...string) *SelectExpect {
	e := SelectExpect(options)

	return &e
}

// MultiSelectExpect expects a multiselect list from console.
type MultiSelectExpect []string

// Do runs the step.
func (e *MultiSelectExpect) Do(c Console) error {
	for _, o := range *e {
		if _, err := c.ExpectString(o); err != nil {
			return err
		}
	}

	return nil
}

// String represents the answer as a string.
func (e *MultiSelectExpect) String() string {
	var sb stringsBuilder

	sb.WriteLinef("Expect a multiselect list:")
	writeSelectList(&sb, *e, multiselectIndicatorRegex)

	return sb.String()
}

func expectMultiSelect(options ...string) *MultiSelectExpect {
	e := MultiSelectExpect(options)

	return &e
}

func breakdownOptions(options []string, indicator *regexp.Regexp) ([]map[string]string, string) {
	breakdown := make([]map[string]string, 0, len(options))

	var size int

	for _, o := range options {
		e := map[string]string{
			"prefix": "",
			"option": "",
		}

		if m := indicator.FindStringSubmatch(o); m != nil {
			e["prefix"] = m[1]
			e["option"] = m[2]
			l := len(m[1])

			if l > size {
				size = l
			}
		} else {
			e["option"] = o
		}

		breakdown = append(breakdown, e)
	}

	return breakdown, fmt.Sprintf("%%-%ds", size)
}

func writeSelectList(sb *stringsBuilder, options []string, indicator *regexp.Regexp) {
	breakdown, pad := breakdownOptions(options, indicator)

	for i, o := range breakdown {
		if i > 0 {
			sb.WriteRune('\n')
		}

		sb.Writef(pad, o["prefix"]).
			Writef(o["option"])
	}
}

// StringExpect expects a string from console.
type StringExpect string

// Do runs the step.
func (e StringExpect) Do(c Console) error {
	_, err := c.ExpectString(string(e))

	return err
}

// String represents the answer as a string.
func (e StringExpect) String() string {
	return fmt.Sprintf("Expect a string: %q", string(e))
}

func expectString(s string) StringExpect {
	return StringExpect(s)
}
