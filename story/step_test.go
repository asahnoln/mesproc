package story_test

import (
	"testing"

	"github.com/asahnoln/mesproc/story"
	"github.com/asahnoln/mesproc/test"
)

func TestOneStep(t *testing.T) {
	stp := story.NewStep().
		Expect("sector 1").
		Respond("The story of this sector").
		Fail("Please type `sector 1`")
	test.AssertSameString(t, "The story of this sector", stp.Response(), "want response %q, got %q")
	test.AssertSameString(t, "Please type `sector 1`", stp.FailMessage(), "want response %q, got %q")
}

func TestExpectGeoLocation(t *testing.T) {
	stp := story.NewStep().
		ExpectGeo(43.257169, 76.924515, 50).
		Respond("Correct location").
		Fail("Location incorrect")

	str := story.New().Add(stp)

	test.AssertSameString(t, stp.Response(), str.RespondWithStepTo(0, "43.257169,76.924515").Text(), "want exact geo response %q, got %q")
	test.AssertSameString(t, stp.Response(), str.RespondWithStepTo(1, "43.257081,76.924835").Text(), "want approximate (50m) geo response %q, got %q")
	test.AssertSameString(t, stp.FailMessage(), str.RespondTo("43.257248572900004,76.92567261243957"), "want fail geo response when far %q, got %q")
}
