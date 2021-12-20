package mesproc_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/asahnoln/mesproc"
)

const command = "/start"

func TestHandle(t *testing.T) {
	h := prepareHandler()
	externalServiceMock(h)

	srv := mesproc.NewServer(h)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"message": "`+command+`"}`))
	srv.ServeHTTP(w, r)

	assertSameString(t, "ok", w.Body.String(), "want response answer %q, got %q")
	assertSameString(t, h.m[command], h.got, "want service receive message %q, got %q")
}

func assertSameString(t testing.TB, want, got, message string) {
	t.Helper()

	if want != got {
		t.Errorf(message, want, got)
	}
}

func assertSameInt(t testing.TB, want, got int, message string) {
	t.Helper()

	if want != got {
		t.Errorf(message, want, got)
	}
}

func prepareHandler() *stubHandler {
	return &stubHandler{m: map[string]string{
		command: "Choose your language",
	}}
}

func externalServiceMock(h *stubHandler) {
	service := httptest.NewServer(h)
	h.target = service.URL
}

type stubHandler struct {
	m      map[string]string
	target string
	got    string
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

func (h *stubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var m struct {
		Message string
	}
	_ = json.NewDecoder(r.Body).Decode(&m)
	h.got = m.Message
}
