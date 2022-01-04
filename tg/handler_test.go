package tg_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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
		step           *story.Step
		update         tg.Update
		tgServerTarget string
	}{
		{story.NewStep().Expect("want this").Respond("standard respond"),
			tg.Update{
				Message: tg.Message{
					Chat: tg.Chat{
						ID: 101,
					},
					Text: "want this",
				},
			}, "/sendMessage"},
		{story.NewStep().Expect("want audio").Respond("audio:http://example.com/audio.mp3"),
			tg.Update{
				Message: tg.Message{
					Chat: tg.Chat{
						ID: 187,
					},
					Text: "want audio",
				},
			}, "/sendAudio"},
		{story.NewStep().ExpectGeo(43, 75, 0).Respond("good").Fail("not good"),
			tg.Update{
				Message: tg.Message{
					Chat: tg.Chat{
						ID: 11,
					},
					Location: &tg.Location{
						Latitude:  43,
						Longitude: 75,
					},
				},
			}, "/sendMessage",
		},
	}

	stg := &stubTgServer{}
	close, target := stg.tgServerMockURL()
	defer close()

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%q to %q", tt.step.Response(), tt.tgServerTarget), func(t *testing.T) {
			str := story.New().Add(tt.step)
			th := tg.New(target, str)

			body, _ := json.Marshal(&tt.update)

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(body))

			th.ServeHTTP(w, r)

			test.AssertSameString(t, tt.tgServerTarget, stg.gotPath, "want tg service called path %q, got %q")
			test.AssertSameString(t, "application/json", stg.gotHeader, "want tg service receiving message %q, got %q")
			test.AssertSameString(t, strings.TrimPrefix(tt.step.Response(), tg.PREFIX_AUDIO), stg.gotText, "want tg service receiving message %q, got %q")
			test.AssertSameInt(t, tt.update.Message.Chat.ID, stg.gotChatID, "want tg service receiving chat id %v, got %v")
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
