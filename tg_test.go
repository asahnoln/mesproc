package mesproc_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/asahnoln/mesproc"
)

type stubTgServer struct {
	gotText, gotHeader, gotPath string
	gotChatID                   int
}

func TestTgHandler(t *testing.T) {
	stg := &stubTgServer{}
	close, target := stg.tgServerMockURL()
	defer close()

	str := mesproc.NewStory().Add(
		mesproc.NewStep().Expect("want this").Respond("Ok you can want it"),
	)
	tg := mesproc.NewTgHandler(target, str)
	srv := mesproc.NewServer(tg)

	update := mesproc.TgUpdate{
		Message: mesproc.TgMessage{
			Chat: mesproc.TgChat{
				ID: 187,
			},
			Text: str.Step().Expectation(),
		},
	}
	body, _ := json.Marshal(&update)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(body))

	// TODO: This returns different results depending on running Story.RespondTo
	want := str.Step().Response()

	srv.ServeHTTP(w, r)

	assertSameString(t, "/sendMessage", stg.gotPath, "want tg service called path %q, got %q")
	assertSameString(t, "application/json", stg.gotHeader, "want tg service receiving message %q, got %q")
	assertSameString(t, want, stg.gotText, "want tg service receiving message %q, got %q")
	assertSameInt(t, update.Message.Chat.ID, stg.gotChatID, "want tg service receiving chat id %v, got %v")
}

func TestTgSendAudio(t *testing.T) {
	stg := &stubTgServer{}
	close, target := stg.tgServerMockURL()
	defer close()

	want := "http://example.com/audio.mp3"
	str := mesproc.NewStory().Add(
		mesproc.NewStep().Expect("want audio").Respond("audio:" + want),
	)
	tg := mesproc.NewTgHandler(target, str)
	srv := mesproc.NewServer(tg)

	update := mesproc.TgUpdate{
		Message: mesproc.TgMessage{
			Chat: mesproc.TgChat{
				ID: 187,
			},
			Text: str.Step().Expectation(),
		},
	}
	body, _ := json.Marshal(&update)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(body))

	srv.ServeHTTP(w, r)

	assertSameString(t, "/sendAudio", stg.gotPath, "want tg service called path %q, got %q")
	assertSameString(t, want, stg.gotText, "want tg service receiving audio link %q, got %q")
}

func (s *stubTgServer) tgServerMockURL() (func(), string) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux := http.NewServeMux()
		s.gotPath = r.URL.Path
		mux.HandleFunc("/sendMessage", func(w http.ResponseWriter, r *http.Request) {
			var m mesproc.TgSendMessage
			json.NewDecoder(r.Body).Decode(&m)
			s.gotHeader = r.Header.Get("Content-Type")
			s.gotText = m.Text
			s.gotChatID = m.ChatID
		})
		mux.HandleFunc("/sendAudio", func(w http.ResponseWriter, r *http.Request) {
			var m mesproc.TgSendAudio
			json.NewDecoder(r.Body).Decode(&m)
			s.gotHeader = r.Header.Get("Content-Type")
			s.gotText = m.Audio
			s.gotChatID = m.ChatID
		})

		mux.ServeHTTP(w, r)
	}))

	return srv.Close, srv.URL
}
