package mesproc_test

import (
	"testing"

	"github.com/asahnoln/mesproc"
)

func TestRespond(t *testing.T) {
	m := mesproc.AnswerMap{
		"/start": "Choose language",
		"/ru":    "Choose sector",
	}
	s := mesproc.NewStory(m)

	got := s.Respond("/start")
	assertSameString(t, "Choose language", got, "want response %q, got %q")

	got = s.Respond("/ru")
	assertSameString(t, "Choose sector", got, "want response %q, got %q")
}
