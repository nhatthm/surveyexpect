package surveymock_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/surveymock"
)

func TestSurvey_ResetExpectations(t *testing.T) {
	t.Parallel()

	s := surveymock.Mock(func(s *surveymock.Survey) {
		s.ExpectPassword("Enter your password:").Times(3)
	})(T())

	s.ResetExpectations()

	assert.NoError(t, s.ExpectationsWereMet())
}
