package tg_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/asahnoln/mesproc/pkg/story"
	"github.com/asahnoln/mesproc/pkg/tg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stubTgServer struct {
	gotText, gotHeader, gotPath []string
	gotChatID                   []int
}

func TestHandler(t *testing.T) {
	tests := []struct {
		step   *story.Step
		update tg.Update
	}{
		{story.NewStep().Expect("want this").Respond("standard respond"),
			tg.Update{
				Message: tg.Message{
					Chat: tg.Chat{
						ID: 101,
					},
					Text: "want this",
				},
			},
		},
		{story.NewStep().Expect("want multi").Respond("first", "second"),
			tg.Update{
				Message: tg.Message{
					Chat: tg.Chat{
						ID: 101,
					},
					Text: "want multi",
				},
			},
		},
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
			},
		},
	}

	for _, tt := range tests {
		stg := &stubTgServer{}
		close, target := stg.tgServerMockURL()
		defer close()

		t.Run(tt.step.Response(), func(t *testing.T) {
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
				assert.Equal(t, "/sendMessage", stg.gotPath[i], "want tg service right path")
				assert.Equal(t, "application/json", stg.gotHeader[i], "want tg service right header")
				assert.Equal(t, r, stg.gotText[i], "want tg service receiving right message")
				assert.Equal(t, tt.update.Message.Chat.ID, stg.gotChatID[i], "want tg service receiving right chat")
			}
		})

	}
}

func TestAudio(t *testing.T) {
	str := story.New().Add(story.NewStep().Respond("audio:http://example.com/audio.mp3").Expect("music").Fail("wrong"))
	stg := &stubTgServer{}
	close, target := stg.tgServerMockURL()
	defer close()

	th := tg.New(target, str, nil)

	body, _ := json.Marshal(tg.Update{
		Message: tg.Message{
			Chat: tg.Chat{
				ID: 6,
			},
			Text: "music",
		},
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	th.ServeHTTP(w, r)

	assert.Equal(t, "/sendChatAction", stg.gotPath[0], "want first to be sent - chat action")
	assert.Equal(t, "application/json", stg.gotHeader[0])
	assert.Equal(t, 6, stg.gotChatID[0])
	assert.Equal(t, "upload_document", stg.gotText[0])

	assert.Equal(t, "/sendAudio", stg.gotPath[1], "want second to be sent - audio")
	assert.Equal(t, "application/json", stg.gotHeader[1])
	assert.Equal(t, 6, stg.gotChatID[1])
	assert.Equal(t, "http://example.com/audio.mp3", stg.gotText[1])
}

func TestPhoto(t *testing.T) {
	str := story.New().Add(story.NewStep().Respond("photo:http://example.com/photo.jpg").Expect("picture").Fail("wrong picture"))
	stg := &stubTgServer{}
	close, target := stg.tgServerMockURL()
	defer close()

	th := tg.New(target, str, nil)

	body, _ := json.Marshal(tg.Update{
		Message: tg.Message{
			Chat: tg.Chat{
				ID: 7,
			},
			Text: "picture",
		},
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	th.ServeHTTP(w, r)

	assert.Equal(t, "/sendChatAction", stg.gotPath[0], "want first to be sent - chat action")
	assert.Equal(t, "application/json", stg.gotHeader[0])
	assert.Equal(t, 7, stg.gotChatID[0])
	assert.Equal(t, "upload_photo", stg.gotText[0])

	t.Logf("%#v\n", stg)
	assert.Equal(t, "/sendPhoto", stg.gotPath[1], "want second to be sent - photo")
	assert.Equal(t, "application/json", stg.gotHeader[1])
	assert.Equal(t, 7, stg.gotChatID[1])
	assert.Equal(t, "http://example.com/photo.jpg", stg.gotText[1])
}

func TestDifferentUsersStepsAndLangs(t *testing.T) {
	str := story.New().
		AddCommand(story.NewStep().Expect("start").Respond("startCommand").Fail("no fail")).
		Add(story.NewStep().Expect("step 1").Respond("you're great!", "go to step 2").Fail("still step 1")).
		Add(story.NewStep().Expect("step 2").Respond("go to step 3").Fail("still step 2")).
		Add(story.NewStep().Expect("step 3").Respond("finish", "non loc finish", "loc finish").Fail("still step 3")).
		I18n(story.I18nMap{
			"ru": {
				"startCommand":            "стартоваяКоманда",
				"step 1":                  "шаг 1",
				"you're great!":           "вы классный!",
				"go to step 2":            "идите к шагу 2",
				"still step 1":            "все еще шаг 1",
				"step 2":                  "шаг 2",
				"still step 2":            "все еще шаг 2",
				"finish":                  "финиш",
				"loc finish":              "loc финиш",
				story.I18nLanguageChanged: "Язык изменен",
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
				From: tg.From{
					LanguageCode: "ru",
				},
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

	tests := []struct {
		id        int
		message   string
		responses []string
	}{
		{2, "/start", []string{"стартоваяКоманда"}},
		{2, "/en", []string{"startCommand"}},

		// Tested different steps for users
		{1, "/en", []string{"Language changed"}}, // TODO: Is it ok it returns this if it was first command ever?
		{1, "step 1", []string{"you're great!", "go to step 2"}},
		{2, "wrong step", []string{"still step 1"}},

		// Testing different languages
		{1, "step 2", []string{"go to step 3"}},

		// Previous response translated
		{2, "/ru", []string{"все еще шаг 1"}},
		{2, "неверно", []string{"все еще шаг 1"}},

		{1, "where am I", []string{"still step 3"}},

		{2, "шаг 1", []string{"вы классный!", "идите к шагу 2"}},
		{2, "/en", []string{"you're great!", "go to step 2"}},
		{2, "/ru", []string{"вы классный!", "идите к шагу 2"}},
		{2, "шаг 2", []string{"go to step 3"}},
		{2, "step 3", []string{"финиш", "non loc finish", "loc финиш"}},
		{2, "что", []string{"все еще шаг 1"}},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("User %d: %q", tt.id, tt.message), func(t *testing.T) {
			sendAndAssert(t, tt.id, tt.message, tt.responses...)
		})
	}
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

func TestLanguageFromTgCode(t *testing.T) {
	stg := &stubTgServer{}
	close, target := stg.tgServerMockURL()
	defer close()

	str := story.New().
		Add(story.NewStep().Expect("Hi!").Respond("nice").Fail("wrong")).
		I18n(story.I18nMap{
			"ru": {
				"nice": "отлично",
			},
		})

	obj := tg.Update{
		Message: tg.Message{
			Chat: tg.Chat{
				ID: 657,
			},
			Text: "Hi!",
			From: tg.From{
				LanguageCode: "ru",
			},
		},
	}
	body, _ := json.Marshal(obj)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))

	th := tg.New(target, str, nil)
	th.ServeHTTP(w, r)

	assert.Equal(t, "отлично", stg.gotText[0], "want immediately russian text because of language code")
}

func TestLaterMessage(t *testing.T) {
	stg := &stubTgServer{}
	close, target := stg.tgServerMockURL()
	defer close()

	str := story.New().
		Add(story.NewStep().
			Expect("want late").
			Respond("first", "second", "late", "immediately").
			Fail("fail").
			Additional(2, "time", time.Millisecond*100)).
		Add(story.NewStep().
			Expect("no expectation").
			Respond("unreachable!").
			Fail("should be unreachable"))

	th := tg.New(target, str, nil)

	obj := tg.Update{
		Message: tg.Message{
			Chat: tg.Chat{
				ID: 657,
			},
			Text: "want late",
		},
	}
	body, _ := json.Marshal(obj)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	th.ServeHTTP(w, r)

	assert.Len(t, stg.gotText, 3)

	// TODO: Should not wait for real
	time.Sleep(time.Millisecond * 202)
	assert.Len(t, stg.gotText, 4)
}

func TestCancelLaterMessage(t *testing.T) {
	stg := &stubTgServer{}
	close, target := stg.tgServerMockURL()
	defer close()

	str := story.New().
		Add(story.NewStep().
			Expect("want late cancel").
			Respond("first", "second", "late but cancelled", "even later", "immediately").
			Fail("fail").
			Additional(2, "time", time.Millisecond*100).
			Additional(3, "time", time.Millisecond*200)).
		Add(story.NewStep().
			Expect("no expectation").
			Respond("unreachable!").
			Fail("should be unreachable"))

	th := tg.New(target, str, nil)

	obj := tg.Update{
		Message: tg.Message{
			Chat: tg.Chat{
				ID: 657,
			},
			Text: "want late cancel",
		},
	}
	body, _ := json.Marshal(obj)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	th.ServeHTTP(w, r)

	assert.Len(t, stg.gotText, 3)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	th.ServeHTTP(w, r)
	time.Sleep(time.Millisecond * 1)

	require.Len(t, stg.gotText, 4)
	assert.Equal(t, "late but cancelled", stg.gotText[3])

	// TODO: Should not wait for real
	time.Sleep(time.Millisecond * 202)
	require.Len(t, stg.gotText, 5)
	assert.Equal(t, "even later", stg.gotText[4])

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	th.ServeHTTP(w, r)
	require.Len(t, stg.gotText, 6)
	assert.Equal(t, "should be unreachable", stg.gotText[5])
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
		// TODO: When testing new uris easy to forget to add them, should be error or something?
		mux.HandleFunc("/sendPhoto", func(w http.ResponseWriter, r *http.Request) {
			var m tg.SendPhoto
			json.NewDecoder(r.Body).Decode(&m)
			fillData(m.ChatID, m.Photo, r)
		})
		mux.HandleFunc("/sendChatAction", func(w http.ResponseWriter, r *http.Request) {
			var m tg.SendChatAction
			json.NewDecoder(r.Body).Decode(&m)
			fillData(m.ChatID, m.Action, r)
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
