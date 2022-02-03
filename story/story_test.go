package story_test

import (
	"testing"

	"github.com/asahnoln/mesproc/story"
	"github.com/asahnoln/mesproc/test"
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
		test.AssertSameString(t, "finish", str.RespondWithStepTo(1, "step 2"), "want response %q, got %q")

		// Assert that function didn't advacne current step
		test.AssertSameString(t, "step 1", str.Step().Expectation(), "want current step response %q, got %q")
	})

}
