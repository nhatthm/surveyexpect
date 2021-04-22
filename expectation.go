package surveyexpect

import (
	"fmt"
	"regexp"
)

var indicatorRegex = regexp.MustCompile(`^([^ ]\s+)(.*)`)

// SelectExpect expects strings from console.
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
	breakdown := make([]map[string]string, 0)

	var size int

	for _, o := range *e {
		e := map[string]string{
			"prefix": "",
			"option": "",
		}

		if m := indicatorRegex.FindStringSubmatch(o); m != nil {
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

	var (
		sb  stringsBuilder
		pad = fmt.Sprintf("%%-%ds", size)
	)

	sb.WriteLinef("Expect a select list:")

	for i, o := range breakdown {
		if i > 0 {
			sb.WriteRune('\n')
		}

		sb.Writef(pad, o["prefix"]).
			Writef(o["option"])
	}

	return sb.String()
}

func expectSelect(options ...string) *SelectExpect {
	e := SelectExpect(options)

	return &e
}
