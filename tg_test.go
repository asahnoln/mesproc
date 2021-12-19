package mesproc_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asahnoln/mesproc"
)

type stubTgServer struct {
	got string
}

func TestTgHandler(t *testing.T) {
	// TODO: Use story module
	stg := &stubTgServer{}
	close, target := stg.tgServerMockURL()
	defer close()

	tg := mesproc.NewTgHandler(target)
	srv := mesproc.NewServer(tg)

	tests := []struct {
		updateMessage string
		wantAnswer    string
	}{
		{"/ru", "Выберите сектор"},
		{"/en", "Choose sector"},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("Send update %q, expect answer %q", tt.updateMessage, tt.wantAnswer), func(t *testing.T) {
			update := mesproc.TgUpdate{
				Message: tt.updateMessage,
			}
			body, _ := json.Marshal(&update)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(body))
			srv.ServeHTTP(w, r)

			assertSameString(t, tt.wantAnswer, stg.got, "want tg service receiving message %q, got %q")
		})
	}
}

func (s *stubTgServer) tgServerMockURL() (func(), string) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux := http.NewServeMux()
		mux.HandleFunc("/sendMessage", func(w http.ResponseWriter, r *http.Request) {
			var m mesproc.TgSendMessage
			json.NewDecoder(r.Body).Decode(&m)
			s.got = m.Text
		})

		mux.ServeHTTP(w, r)
	}))

	return srv.Close, srv.URL
}
