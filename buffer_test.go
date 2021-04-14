package surveyexpect

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuffer_Reset(t *testing.T) {
	t.Parallel()

	buf := new(Buffer)

	_, err := buf.Write([]byte("hello world"))
	require.NoError(t, err)

	assert.Equal(t, "hello world", buf.String())

	buf.Reset()
	assert.Empty(t, buf.String())
}
