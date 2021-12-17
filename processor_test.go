package mesproc_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/asahnoln/mesproc"
)

type stubHandler struct {
	m      map[string]string
	target string
}

func (h *stubHandler) Receive(w http.ResponseWriter, r *http.Request) string {
	var m struct {
		Message string
	}
	_ = json.NewDecoder(r.Body).Decode(&m)
	_, _ = w.Write([]byte("ok"))
	return m.Message
}

func (h *stubHandler) Send(k string) {
	_, _ = http.Post(h.target, "", strings.NewReader(`{"message": "`+h.m[k]+`"}`))
}

func TestReceive(t *testing.T) {
	h := &stubHandler{}
	srv := mesproc.NewServer(h)
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
	m := map[string]string{
		"/start": "Choose your language",
	}

	var got string
	service := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m struct {
			Message string
		}
		_ = json.NewDecoder(r.Body).Decode(&m)
		got = m.Message
	}))

	h := &stubHandler{m: m, target: service.URL}
	srv := mesproc.NewServer(h)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"message": "/start"}`))
	srv.ServeHTTP(w, r)

	want := m["/start"]
	if want != got {
		t.Errorf("want service receive message %q, got %q", want, got)
	}
}
