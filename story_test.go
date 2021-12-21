package mesproc_test

import (
	"testing"

	"github.com/asahnoln/mesproc"
)

// TODO: Test finished steps. What happens? Loop?

func TestOneStep(t *testing.T) {
	stp := mesproc.NewStep().
		Expect("sector 1").
		Respond("The story of this sector").
		Fail("Please type `sector 1`")
	assertSameString(t, "The story of this sector", stp.Response(), "want response %q, got %q")
	assertSameString(t, "Please type `sector 1`", stp.FailMessage(), "want response %q, got %q")
}

func TestStoryCurrentStep(t *testing.T) {
	str := mesproc.NewStory().Add(mesproc.NewStep().Expect("wow").Respond("yes!"))

	assertSameString(t, "wow", str.Step().Expectation(), "want current step expectation %q, got %q")
	assertSameString(t, "yes!", str.Step().Response(), "want current step expectation %q, got %q")
}

func TestSteps(t *testing.T) {
	str := mesproc.NewStory()

	stp1 := mesproc.NewStep().
		Expect("sector 1").
		Respond("The story of this sector").
		Fail("Please type `sector 1`")
	stp2 := mesproc.NewStep().
		Expect("lulz").Respond("Chilling lulz").Fail("I want to hear `lulz`")

	str.Add(stp1).Add(stp2)

	// TODO: RespondTo has side effects - winding the current step. Should rethink design?
	assertSameString(t, stp1.FailMessage(), str.RespondTo("smth else"), "want response %q, got %q")
	assertSameString(t, stp1.Response(), str.RespondTo("sector 1"), "want response %q, got %q")

	assertSameString(t, stp2.FailMessage(), str.RespondTo("sector 1"), "want response %q, got %q")
	assertSameString(t, stp2.Response(), str.RespondTo("lulz"), "want response %q, got %q")
}
