package tg_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asahnoln/mesproc/story"
	"github.com/asahnoln/mesproc/test"
	"github.com/asahnoln/mesproc/tg"
)

type stubTgServer struct {
	gotText, gotHeader, gotPath string
	gotChatID                   int
}

func TestHandler(t *testing.T) {
	tests := []struct {
		responsePrefix string
		want           string
		tgServerTarget string
	}{
		{"", "standard response", "/sendMessage"},
		{"audio:", "http://example.com/audio.mp3", "/sendAudio"},
	}

	stg := &stubTgServer{}
	close, target := stg.tgServerMockURL()
	defer close()

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q to %q", tt.want, tt.tgServerTarget), func(t *testing.T) {
			str := story.New().Add(
				story.NewStep().Expect("want this").Respond(tt.responsePrefix + tt.want),
			)
			th := tg.New(target, str)

			update := tg.Update{
				Message: tg.Message{
					Chat: tg.Chat{
						ID: 187,
					},
					Text: str.Step().Expectation(),
				},
			}
			body, _ := json.Marshal(&update)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(body))

			th.ServeHTTP(w, r)

			test.AssertSameString(t, tt.tgServerTarget, stg.gotPath, "want tg service called path %q, got %q")
			test.AssertSameString(t, "application/json", stg.gotHeader, "want tg service receiving message %q, got %q")
			test.AssertSameString(t, tt.want, stg.gotText, "want tg service receiving message %q, got %q")
			test.AssertSameInt(t, update.Message.Chat.ID, stg.gotChatID, "want tg service receiving chat id %v, got %v")
		})

	}
}

func (s *stubTgServer) tgServerMockURL() (func(), string) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux := http.NewServeMux()
		s.gotPath = r.URL.Path
		fillData := func(id int, text string, r *http.Request) {
			s.gotHeader = r.Header.Get("Content-Type")
			s.gotChatID = id
			s.gotText = text
		}
		mux.HandleFunc("/sendMessage", func(w http.ResponseWriter, r *http.Request) {
			var m tg.SendMessage
			json.NewDecoder(r.Body).Decode(&m)
			fillData(m.ChatID, m.Text, r)
		})
		mux.HandleFunc("/sendAudio", func(w http.ResponseWriter, r *http.Request) {
			var m tg.SendAudio
			json.NewDecoder(r.Body).Decode(&m)
			fillData(m.ChatID, m.Audio, r)
		})

		mux.ServeHTTP(w, r)
	}))

	return srv.Close, srv.URL
}
