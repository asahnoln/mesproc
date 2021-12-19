package mesproc_test

import (
	"testing"

	"github.com/asahnoln/mesproc"
)

func TestRespond(t *testing.T) {
	const (
		s1   = "sector 1"
		lulz = "lulz"
		w5   = "winners 5"
	)
	m := mesproc.AnswerMap{
		s1:   "That's the story of sector 1. Now type `lulz`",
		lulz: "Middle story of lulz",
		w5:   "Final story of winners 5",
	}
	s := mesproc.NewStory(m)

	assertSameString(t, m[s1], s.Respond(s1), "want response %q, got %q")
	assertSameString(t, "Please type `lulz`", s.Respond(s1), "want response %q to repeated request, got %q")

	assertSameString(t, m[lulz], s.Respond(lulz), "want response %q, got %q")
	assertSameString(t, "Please type `winners 5`", s.Respond(lulz), "want response %q to repeated request, got %q")

	assertSameString(t, m[w5], s.Respond(w5), "want response %q, got %q")
	assertSameString(t, "Please type `guds`", s.Respond(w5), "want response %q to repeated request, got %q")
}
