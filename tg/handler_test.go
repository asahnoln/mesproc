package tg_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/asahnoln/mesproc/story"
	"github.com/asahnoln/mesproc/tg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubTgServer struct {
	gotText, gotHeader, gotPath []string
	gotChatID                   []int
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
		{story.NewStep().Expect("want multi").Respond("first", "second"),
			tg.Update{
				Message: tg.Message{
					Chat: tg.Chat{
						ID: 101,
					},
					Text: "want multi",
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

	for _, tt := range tests {
		stg := &stubTgServer{}
		close, target := stg.tgServerMockURL()
		defer close()

		t.Run(fmt.Sprintf("%q to %q", tt.step.Responses(), tt.tgServerTarget), func(t *testing.T) {
			str := story.New().Add(tt.step)
			th := tg.New(target, str, nil)

			body, err := json.Marshal(&tt.update)
			require.NoError(t, err, "unexpected error while marshalling object")

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(body))

			th.ServeHTTP(w, r)

			rs := tt.step.Responses()
			require.Len(t, stg.gotText, len(rs), "want the same count of requests to tg server as responses")

			for i, r := range rs {
				assert.Equal(t, tt.tgServerTarget, stg.gotPath[i], "want tg service right path")
				assert.Equal(t, "application/json", stg.gotHeader[i], "want tg service right header")
				assert.Equal(t, strings.TrimPrefix(r, tg.PrefixAudio), stg.gotText[i], "want tg service receiving right message")
				assert.Equal(t, tt.update.Message.Chat.ID, stg.gotChatID[i], "want tg service receiving right chat")
			}
		})

	}
}

func TestDifferentUsersStepsAndLangs(t *testing.T) {
	str := story.New().
		Add(story.NewStep().Expect("step 1").Respond("you're great!", "go to step 2").Fail("still step 1")).
		Add(story.NewStep().Expect("step 2").Respond("go to step 3").Fail("still step 2")).
		Add(story.NewStep().Expect("step 3").Respond("finish", "non loc finish", "loc finish").Fail("still step 3")).
		I18n(story.I18nMap{
			"ru": {
				"step 1":        "шаг 1",
				"you're great!": "вы классный!",
				"go to step 2":  "идите к шагу 2",
				"still step 1":  "все еще шаг 1",
				"step 2":        "шаг 2",
				"still step 2":  "все еще шаг 2",
				"finish":        "финиш",
				"loc finish":    "loc финиш",
			},
		})

	stg := &stubTgServer{}
	close, target := stg.tgServerMockURL()
	defer close()

	th := tg.New(target, str, nil)

	// TODO: Rework through table tests so we don't have this difficult logic of testing
	sendAndAssert := func(t testing.TB, id int, text string, want ...string) {
		t.Helper()
		body, err := json.Marshal(tg.Update{
			Message: tg.Message{
				Chat: tg.Chat{
					ID: id,
				},
				Text: text,
			},
		})
		require.NoError(t, err, "unexpected error while marshaling object")

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
		th.ServeHTTP(w, r)

		for j, w := range want {
			assert.Equal(t, w, stg.gotText[j], "want response for user %d", id)
		}

		// Reset server values
		stg.zero()
	}

	// Tested different steps for users
	sendAndAssert(t, 1, "step 1", "you're great!", "go to step 2")
	sendAndAssert(t, 2, "wrong step", "still step 1")

	// Testing different languages
	sendAndAssert(t, 1, "step 2", "go to step 3")

	sendAndAssert(t, 2, "/ru", "Language changed")
	sendAndAssert(t, 2, "неверно", "все еще шаг 1")

	sendAndAssert(t, 1, "where am I", "still step 3")

	sendAndAssert(t, 2, "шаг 1", "вы классный!", "идите к шагу 2")
	sendAndAssert(t, 2, "шаг 2", "go to step 3")
	sendAndAssert(t, 2, "step 3", "финиш", "non loc finish", "loc финиш")
	sendAndAssert(t, 2, "что", "все еще шаг 1")
}

func TestLogging(t *testing.T) {
	stg := stubTgServer{}
	close, target := stg.tgServerMockURL()
	defer close()

	b := bytes.Buffer{}
	lgr := log.New(&b, "", 0)
	th := tg.New(
		target,
		story.New().Add(story.NewStep().Expect("ok").Respond("what")),
		lgr,
	)

	obj := tg.Update{
		Message: tg.Message{
			Chat: tg.Chat{
				ID: 54,
			},
			Text: "something",
		},
	}
	body, _ := json.Marshal(obj)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
	th.ServeHTTP(w, r)

	assert.Contains(t, b.String(), "telegram update: ", "want logged message on receiving")
	assert.Contains(t, b.String(), fmt.Sprintf("%#v", obj), "want logged update object on receiving")
	assert.Contains(t, b.String(), time.Now().Format(time.RFC3339), "want logged date on receiving")
}

func (s *stubTgServer) tgServerMockURL() (func(), string) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mux := http.NewServeMux()
		s.gotPath = append(s.gotPath, r.URL.Path)
		fillData := func(id int, text string, r *http.Request) {
			s.gotHeader = append(s.gotHeader, r.Header.Get("Content-Type"))
			s.gotChatID = append(s.gotChatID, id)
			s.gotText = append(s.gotText, text)
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

func (s *stubTgServer) zero() {
	s.gotChatID = []int{}
	s.gotHeader = []string{}
	s.gotPath = []string{}
	s.gotText = []string{}
}
