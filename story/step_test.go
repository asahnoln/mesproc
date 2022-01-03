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
