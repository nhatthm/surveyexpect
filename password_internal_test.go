package surveymock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword_Once(t *testing.T) {
	t.Parallel()

	p := newPassword(&Survey{}, "").Once()

	assert.Equal(t, 1, p.repeatability)
}

func TestPassword_Twice(t *testing.T) {
	t.Parallel()

	p := newPassword(&Survey{}, "").Twice()

	assert.Equal(t, 2, p.repeatability)
}

func TestPassword_Times(t *testing.T) {
	t.Parallel()

	p := newPassword(&Survey{}, "").Times(5)

	assert.Equal(t, 5, p.repeatability)
}
