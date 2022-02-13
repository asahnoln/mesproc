package story_test

import (
	"testing"

	"github.com/asahnoln/mesproc/story"
	"github.com/stretchr/testify/assert"
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

// func TestSaveExpectation(t *testing.T) {
// 	store := &stubStore{}
// 	stp := story.NewStep().ExpectSave(store).Respond("thank you!")

// 	str := story.New().Add(stp)

// 	assert.Equal(t, stp.Response(), str.RespondWithLangStepTo(0, "", "I save this"), "want response on saving message")
// 	assert.Equal(t, "I save this", store.message, "want response on saving message")
// }
