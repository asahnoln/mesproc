package story_test

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/asahnoln/mesproc/pkg/story"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSteps(t *testing.T) {
	str := story.New()

	stp1 := story.NewStep().
		Expect("sector 1").
		Respond("The story of this sector").
		Fail("Please type `sector 1`")
	stp2 := story.NewStep().
		Expect("lulz").Respond("Chilling lulz").Fail("I want to hear `lulz`")

	str.Add(stp1).Add(stp2)

	// TODO: RespondTo has side effects - winding the current step. Should rethink design?
	assert.Equal(t, stp1.FailMessage(), str.ResponsesWithLangStepTo(0, "", "smth else")[0].Text(), "want fail response")
	assert.Equal(t, stp1.Response(), str.ResponsesWithLangStepTo(0, "", "sector 1")[0].Text(), "want correct response")

	assert.Equal(t, stp2.FailMessage(), str.ResponsesWithLangStepTo(1, "", "sector 1")[0].Text(), "want response %q, got %q")
	assert.Equal(t, stp2.Response(), str.ResponsesWithLangStepTo(1, "", "lulz")[0].Text(), "want response %q, got %q")
}

func TestStoryI18N(t *testing.T) {
	str := story.New().
		Add(story.NewStep().Expect("somestep").Fail("wrong")).
		I18n(story.I18nMap{
			"ru": {
				story.I18nLanguageChanged: "Язык изменен на русский",
			},
			"kk": {
				story.I18nLanguageChanged: "Язык изменен на казахский",
			},
		})

	assert.Equal(t, "en", str.ResponsesWithLangStepTo(0, "", "say")[0].Lang(), "want English language by default")

	assert.Equal(t, "Язык изменен на русский", str.ResponsesWithLangStepTo(0, "", "/ru")[0].Text(), "want language message in Russian")
	assert.Equal(t, "ru", str.ResponsesWithLangStepTo(0, "en", "/ru")[0].Lang(), "want Russian language")

	assert.Equal(t, "Язык изменен на казахский", str.ResponsesWithLangStepTo(0, "", "/kk")[0].Text(), "want language message in Kazakh")
	assert.Equal(t, "kk", str.ResponsesWithLangStepTo(0, "", "/kk")[0].Lang(), "want Kazakh language")

	assert.Equal(t, story.I18nLanguageChanged, str.ResponsesWithLangStepTo(0, "", "/en")[0].Text(), "want language message in English")
	assert.Equal(t, "en", str.ResponsesWithLangStepTo(0, "", "/en")[0].Lang(), "want English language")
}

func TestSettingLanguageChangesMessage(t *testing.T) {
	str := story.New().
		Add(story.NewStep().Expect("step 1").Respond("story 1").Fail("please write step 1")).
		Add(story.NewStep().Expect("step 2").Respond("story 2").Fail("please write step 2")).
		Add(story.NewStep().Expect("step 3").Respond("story 3").Fail("please write step 3")).
		I18n(story.I18nMap{
			"ru": {
				"step 1":              "шаг 1",
				"story 1":             "история 1",
				"step 2":              "шаг 2",
				"please write step 3": "введите шаг 3",
			},
		})

	assert.Equal(t, "история 1", str.ResponsesWithLangStepTo(0, "ru", "шаг 1")[0].Text(), "want i18n response %q, got %q")

	// Default line should come out if there is no i18n
	assert.Equal(t, "story 2", str.ResponsesWithLangStepTo(1, "ru", "шаг 2")[0].Text(), "want i18n default response %q, got %q")

	assert.Equal(t, "введите шаг 3", str.ResponsesWithLangStepTo(2, "ru", "wrong")[0].Text(), "want i18n fail response %q, got %q")

	assert.Equal(t, "ru", str.ResponsesWithLangStepTo(0, "en", "/ru")[0].Lang(), "want i18n successful change")
}

func TestRespondToSpecificStep(t *testing.T) {
	str := story.New().
		Add(story.NewStep().Expect("step 1").Respond("go to step 2").Fail("still step 1")).
		Add(story.NewStep().Expect("step 2").Respond("finish").Fail("still step 2"))

	t.Run("step rotates if out of range", func(t *testing.T) {
		assert.Equal(t, "go to step 2", str.ResponsesWithLangStepTo(2, "", "step 1")[0].Text(), "want response %q, got %q")
		assert.Equal(t, "still step 2", str.ResponsesWithLangStepTo(3, "", "wrong step")[0].Text(), "want response %q, got %q")
	})

	t.Run("successful response should advance the step", func(t *testing.T) {
		assert.True(t, str.ResponsesWithLangStepTo(4, "", "step 1")[0].ShouldAdvance(), "want ShouldAdvance() = true")
	})

	t.Run("wrong expectation must not advance the step", func(t *testing.T) {
		assert.False(t, str.ResponsesWithLangStepTo(5, "", "step 1")[0].ShouldAdvance(), "want ShouldAdvance() = false")
	})
}

func TestRespondToCommand(t *testing.T) {
	stp := story.NewStep().Expect("command").Respond("that was a command").Fail("failed processing the command")
	str := story.New().AddCommand(stp)
	str.I18n(story.I18nMap{
		"ru": {
			"that was a command": "это была команда",
		},
	})

	assert.Equal(t, stp.Response(), str.ResponsesWithLangStepTo(5, "", "/command")[0].Text(), "want response %q, got %q")

	assert.Equal(t, "это была команда", str.ResponsesWithLangStepTo(6, "ru", "/command")[0].Text(), "want response %q, got %q")
}

func TestLoadingFromJSON(t *testing.T) {
	err := os.Mkdir("testdata/save", 0755)
	require.NoError(t, err, "unexpected error creating test dir")
	defer os.RemoveAll("testdata/save")

	f, err := os.Open("testdata/story.json")
	require.NoError(t, err, "unexpected error loading test file")
	defer f.Close()

	str, err := story.Load(f)
	require.NoError(t, err, "unexpected error when loading proper JSON for the story")

	t.Run("Common steps", func(t *testing.T) {
		assert.Equal(t, "still at step 1", str.ResponsesWithLangStepTo(0, "", "help")[0].Text(), "want fail message in response to wrong expectation")
		assert.Equal(t, "now at step 2", str.ResponsesWithLangStepTo(0, "", "go to step 2")[0].Text(), "want response message to expectation")
		assert.Equal(t, "proper geo", str.ResponsesWithLangStepTo(1, "", "43.257081,76.924835")[0].Text(), "want successfule response to approximate (50m) geo expectation")
		assert.Equal(t, "now finished", str.ResponsesWithLangStepTo(2, "", "finish")[0].Text(), "want response message to final expectation")
	})

	t.Run("Commands", func(t *testing.T) {
		assert.Equal(t, "let's start", str.ResponsesWithLangStepTo(99, "", "/start")[0].Text(), "want response message to command")
	})

	t.Run("Multi response", func(t *testing.T) {
		rs := str.ResponsesWithLangStepTo(3, "", "multi")
		assert.Len(t, rs, 3, "want multi response step")
	})

	t.Run("Saving", func(t *testing.T) {
		assert.Equal(t, "saved!", str.ResponsesWithLangStepTo(4, "", "I want this saved")[0].Text(), "want response message to saving expectation")
	})

	t.Run("Additional info", func(t *testing.T) {
		rs := str.ResponsesWithLangStepTo(3, "", "multi")
		assert.Equal(t, time.Second*600, rs[2].Additional["time"], "want time field on 3rd response of 4th step")
	})
}

func TestErrorLoadingFromJSON(t *testing.T) {
	_, err := story.Load(strings.NewReader(""))

	require.Error(t, err, "want error when loading wrong json")
}

func TestExpectationCaseInsensitive(t *testing.T) {
	str := story.New().Add(story.NewStep().Expect("LOL this GOOD").Respond("success!").Fail("failed"))

	assert.Equal(t, "success!", str.ResponsesWithLangStepTo(0, "", "lOl ThIs GoOd")[0].Text(), "want case-insensitive expectation")
}

func TestSeveralResponses(t *testing.T) {
	str := story.New().
		Add(story.NewStep().Expect("message").Respond("first message", "second message")).
		I18n(story.I18nMap{
			"ru": {
				"message":        "сообщение",
				"first message":  "первое сообщение",
				"second message": "второе сообщение",
			},
		})

	rs := str.ResponsesWithLangStepTo(0, "", "message")
	assert.Len(t, rs, 2, "want several responses")
	assert.Equal(t, "first message", rs[0].Text(), "want first message")
	assert.Equal(t, "second message", rs[1].Text(), "want second message")

	rs = str.ResponsesWithLangStepTo(0, "ru", "сообщение")
	assert.Len(t, rs, 2, "want several responses")
	assert.Equal(t, "первое сообщение", rs[0].Text(), "want first message")
	assert.Equal(t, "второе сообщение", rs[1].Text(), "want second message")
}
