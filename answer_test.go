package surveyexpect

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoAnswer_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "<no answer>", noAnswer().String())
}

func TestInterruptAnswer_String(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "<interrupt>", interruptAnswer().String())
}
