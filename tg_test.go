package mesproc_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asahnoln/mesproc"
)

type stubTgServer struct {
	got string
}

func TestTgHandler(t *testing.T) {
	stg := &stubTgServer{}
	tg := mesproc.NewTgHandler(stg.tgServerMockURL())
	srv := mesproc.NewServer(tg)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	srv.ServeHTTP(w, r)

	assertSameString(t, "Choose sector", stg.got, "want tg service receiving message %q, got %q")
}

func (s *stubTgServer) tgServerMockURL() string {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var m mesproc.TgSendMessage
		json.NewDecoder(r.Body).Decode(&m)
		s.got = m.Text
	}))
	return srv.URL
}
