package mesproc_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/asahnoln/mesproc"
)

func TestReceive(t *testing.T) {
	srv := mesproc.NewServer(nil, "")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"message": "/start"}`))
	srv.ServeHTTP(w, r)

	want := "ok"
	got := w.Body.String()

	if want != got {
		t.Errorf("want response answer %q, got %q", want, got)
	}
}

func TestSend(t *testing.T) {
	m := mesproc.AnswerMap{
		"/start": "Choose your language",
	}
	var got string
	service := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m struct {
			Message string
		}
		json.NewDecoder(r.Body).Decode(&m)
		got = m.Message
	}))

	srv := mesproc.NewServer(m, service.URL)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"message": "/start"}`))
	srv.ServeHTTP(w, r)

	want := m["/start"]
	if want != got {
		t.Errorf("want service receive message %q, got %q", want, got)
	}
}
