package story_test

import (
	"os"
	"strings"
	"testing"

	"github.com/asahnoln/mesproc/story"
	"github.com/asahnoln/mesproc/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoryCurrentStep(t *testing.T) {
	str := story.New().Add(story.NewStep().Expect("wow").Respond("yes!"))

	test.AssertSameString(t, "wow", str.Step().Expectation(), "want current step expectation %q, got %q")
	test.AssertSameString(t, "yes!", str.Step().Response(), "want current step expectation %q, got %q")
}

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
	test.AssertSameString(t, stp1.FailMessage(), str.RespondTo("smth else"), "want response %q, got %q")
	test.AssertSameString(t, stp1.Response(), str.RespondTo("sector 1"), "want response %q, got %q")

	test.AssertSameString(t, stp2.FailMessage(), str.RespondTo("sector 1"), "want response %q, got %q")
	test.AssertSameString(t, stp2.Response(), str.RespondTo("lulz"), "want response %q, got %q")
}

func TestStoryMustLoop(t *testing.T) {
	str := story.New().
		Add(story.NewStep().Expect("s1").Respond("step 1")).
		Add(story.NewStep().Expect("s2").Respond("step 2")).
		Add(story.NewStep().Expect("s3").Respond("step 3"))

	str.RespondTo("s1")
	str.RespondTo("s2")
	str.RespondTo("s3")

	test.AssertSameString(t, "step 1", str.RespondTo("s1"), "want response %q got %q")
}

func TestStoryI18N(t *testing.T) {
	str := story.New().
		I18n(story.I18nMap{
			"ru": {
				story.I18nLanguageChanged: "Язык изменен на русский",
			},
			"kk": {
				story.I18nLanguageChanged: "Язык изменен на казахский",
			},
		})

	test.AssertSameString(t, "Язык изменен на русский", str.RespondTo("/ru"), "want language message %q, got %q")
	test.AssertSameString(t, "Язык изменен на казахский", str.RespondTo("/kk"), "want language message %q, got %q")
	test.AssertSameString(t, story.I18nLanguageChanged, str.RespondTo("/en"), "want language message %q, got %q")
	test.AssertSameString(t, "Язык изменен на русский", str.RespondWithStepTo(12, "/ru").Text(), "want language message %q, got %q")
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

	str.RespondTo("/ru")

	test.AssertSameString(t, "ru", str.Language(), "want language %q, got %q")
	test.AssertSameString(t, "история 1", str.RespondTo("шаг 1"), "want i18n response %q, got %q")

	// Default line should come out if there is no i18n
	test.AssertSameString(t, "story 2", str.RespondTo("шаг 2"), "want i18n default response %q, got %q")

	test.AssertSameString(t, "введите шаг 3", str.RespondTo("wrong"), "want i18n fail response %q, got %q")
}

func TestRespondToSpecificStep(t *testing.T) {
	str := story.New().
		Add(story.NewStep().Expect("step 1").Respond("go to step 2").Fail("still step 1")).
		Add(story.NewStep().Expect("step 2").Respond("finish").Fail("still step 2"))

	t.Run("doesn't affect current step of story", func(t *testing.T) {
		test.AssertSameString(t, "finish", str.RespondWithStepTo(1, "step 2").Text(), "want response %q, got %q")

		// Assert that function didn't advacne current step
		test.AssertSameString(t, "step 1", str.Step().Expectation(), "want current step response %q, got %q")
	})

	t.Run("step rotates if out of range", func(t *testing.T) {
		test.AssertSameString(t, "go to step 2", str.RespondWithStepTo(2, "step 1").Text(), "want response %q, got %q")
		test.AssertSameString(t, "still step 2", str.RespondWithStepTo(3, "wrong step").Text(), "want response %q, got %q")
	})

	t.Run("successful response should advance the step", func(t *testing.T) {
		assert.True(t, str.RespondWithStepTo(4, "step 1").ShouldAdvance(), "want ShouldAdvance() = true")
	})

	t.Run("wrong expectation must not advance the step", func(t *testing.T) {
		assert.False(t, str.RespondWithStepTo(5, "step 1").ShouldAdvance(), "want ShouldAdvance() = false")
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

	test.AssertSameString(t, stp.Response(), str.RespondWithStepTo(5, "/command").Text(), "want response %q, got %q")

	str.SetLanguage("ru")
	test.AssertSameString(t, "это была команда", str.RespondWithStepTo(6, "/command").Text(), "want response %q, got %q")
}

func TestLoadingFromJSON(t *testing.T) {
	f, err := os.Open("testdata/story.json")
	defer f.Close()

	require.NoError(t, err, "error opening file")

	str, err := story.New().Load(f)
	require.NoError(t, err, "unexpected error when loading proper JSON for the story")

	assert.Equal(t, "still at step 1", str.RespondWithStepTo(0, "help").Text(), "want fail message in response to wrong expectation")
	assert.Equal(t, "now at step 2", str.RespondWithStepTo(0, "go to step 2").Text(), "want response message to expectation")
	assert.Equal(t, "proper geo", str.RespondWithStepTo(1, "43.257081,76.924835").Text(), "want successfule response to approximate (50m) geo expectation")
	assert.Equal(t, "now finished", str.RespondWithStepTo(2, "finish").Text(), "want response message to final expectation")
}

func TestErrorLoadingFromJSON(t *testing.T) {
	_, err := story.New().Load(strings.NewReader(""))

	require.Error(t, err, "want error when loading wrong json")
}
