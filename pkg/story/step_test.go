package story_test

import (
	"errors"
	"testing"

	"github.com/asahnoln/mesproc/pkg/story"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOneStep(t *testing.T) {
	stp := story.NewStep().
		Expect("sector 1").
		Respond("The story of this sector").
		Fail("Please type `sector 1`")
	assert.Equal(t, "The story of this sector", stp.Response(), "want proper step response")
	assert.Equal(t, "Please type `sector 1`", stp.FailMessage(), "want proper step fail message")
}

func TestExpectGeoLocation(t *testing.T) {
	stp := story.NewStep().
		ExpectGeo(43.257169, 76.924515, 50).
		Respond("Correct location").
		Fail("Location incorrect")

	str := story.New().Add(stp)

	assert.Equal(t, stp.Response(), str.ResponsesWithLangStepTo(0, "", "43.257169,76.924515")[0].Text(), "want exact geo response")
	assert.Equal(t, stp.Response(), str.ResponsesWithLangStepTo(1, "", "43.257081,76.924835")[0].Text(), "want approximate (50m) geo response")
	assert.Equal(t, stp.FailMessage(), str.ResponsesWithLangStepTo(0, "", "43.257248572900004,76.92567261243957")[0].Text(), "want fail geo response when far")
}

// TODO: Merge with error test
func TestSaveExpectation(t *testing.T) {
	store := &stubStore{}
	store.On("Save", "I save this").Return(nil)

	stp := story.NewStep().ExpectSave(store).Respond("thank you!").Fail("Shouldn't fail but failed")
	str := story.New().Add(stp)

	assert.Equal(t, stp.Response(), str.ResponsesWithLangStepTo(0, "", "I save this")[0].Text(), "want response on saving message")
	store.AssertExpectations(t)
}

// TODO: Need to log error
func TestSaveExpectationError(t *testing.T) {
	store := &stubStore{}
	store.On("Save", mock.Anything).Return(errors.New("save fail"))

	stp := story.NewStep().ExpectSave(store).Respond("thank you!").Fail("Try again")
	str := story.New().Add(stp)

	assert.Equal(t, stp.FailMessage(), str.ResponsesWithLangStepTo(0, "", "Bad will happen")[0].Text(), "want response on saving message")
	store.AssertExpectations(t)
}

type stubStore struct {
	mock.Mock
}

func (s *stubStore) Save(m string) error {
	args := s.Called(m)
	return args.Error(0)
}
