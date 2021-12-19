package mesproc_test

import (
	"testing"

	"github.com/asahnoln/mesproc"
)

func TestSteps(t *testing.T) {
	str := mesproc.NewStory()

	stp1 := mesproc.NewStep().Expect("sector 1").Respond("The story of this sector")
	assertSameString(t, "The story of this sector", stp1.Response(), "want response %q, got %q")

	stp2 := mesproc.NewStep().Expect("lulz").Respond("Chilling lulz")

	str.Add(stp1)
	str.Add(stp2)

	assertSameString(t, stp1.Response(), str.RespondTo("sector 1"), "want response %q, got %q")
	assertSameString(t, stp2.Response(), str.RespondTo("lulz"), "want response %q, got %q")
}
