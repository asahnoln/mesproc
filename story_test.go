package mesproc_test

import (
	"testing"

	"github.com/asahnoln/mesproc"
)

func TestRespond(t *testing.T) {
	const (
		s1 = "sector 1"
		s2 = "sector 2"
	)
	m := mesproc.AnswerMap{
		s1: "That's the story of sector 1. Now type `sector 2`",
		s2: "Final story of sector 2",
	}
	s := mesproc.NewStory(m)

	got := s.Respond(s1)
	assertSameString(t, m[s1], got, "want response %q, got %q")

	got = s.Respond(s2)
	assertSameString(t, m[s2], got, "want response %q, got %q")
}
