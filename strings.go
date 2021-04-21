package surveyexpect

import (
	"fmt"
	"strings"
)

type stringsBuilder struct {
	strings.Builder
}

func (b *stringsBuilder) WriteRune(r rune) *stringsBuilder {
	_, _ = b.Builder.WriteRune(r)

	return b
}

func (b *stringsBuilder) WriteString(s string) *stringsBuilder { // nolint: unparam
	_, _ = b.Builder.WriteString(s)

	return b
}

func (b *stringsBuilder) Writef(format string, args ...interface{}) *stringsBuilder {
	_, _ = fmt.Fprintf(b, format, args...)

	return b
}

func (b *stringsBuilder) WriteLinef(format string, args ...interface{}) *stringsBuilder { // nolint: unparam
	return b.Writef(format, args...).
		WriteRune('\n')
}

func (b *stringsBuilder) WriteLabelLinef(label, value string, args ...interface{}) *stringsBuilder {
	b.Writef("%-7s: ", label).
		WriteLinef(value, args...)

	return b
}
