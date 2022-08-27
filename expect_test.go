package surveyexpect_test

import (
	"fmt"

	"go.nhat.io/surveyexpect"
)

type TestingT struct {
	error *surveyexpect.Buffer
	log   *surveyexpect.Buffer

	clean func()
}

func (t *TestingT) Errorf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.error, format, args...)
}

func (t *TestingT) Log(args ...interface{}) {
	_, _ = fmt.Fprintln(t.log, args...)
}

func (t *TestingT) Logf(format string, args ...interface{}) {
	_, _ = fmt.Fprintf(t.log, format, args...)
}

func (t *TestingT) FailNow() {
	panic("failed")
}

func (t *TestingT) Cleanup(clean func()) {
	t.clean = clean
}

func (t *TestingT) ErrorString() string {
	return t.error.String()
}

func (t *TestingT) LogString() string {
	return t.log.String()
}

func T() *TestingT {
	return &TestingT{
		error: new(surveyexpect.Buffer),
		log:   new(surveyexpect.Buffer),
		clean: func() {},
	}
}
