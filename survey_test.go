package surveyexpect_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nhatthm/surveyexpect"
)

func TestSurvey_ExpectNoExpectation(t *testing.T) {
	t.Parallel()

	s := surveyexpect.Expect()(t)
	err := s.Expect(nil)

	assert.Equal(t, surveyexpect.ErrNoExpectation, err)
}

func TestSurvey_ResetExpectations(t *testing.T) {
	t.Parallel()

	s := surveyexpect.Expect(func(s *surveyexpect.Survey) {
		s.ExpectPassword("Enter your password:").Times(3)
	})(T())

	s.ResetExpectations()

	assert.NoError(t, s.ExpectationsWereMet())
}
