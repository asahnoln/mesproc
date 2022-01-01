package mesproc_test

import (
	"testing"

	"github.com/asahnoln/mesproc"
)

// TODO: Test finished steps. What happens? Loop?

func TestOneStep(t *testing.T) {
	stp := mesproc.NewStep().
		Expect("sector 1").
		Respond("The story of this sector").
		Fail("Please type `sector 1`")
	assertSameString(t, "The story of this sector", stp.Response(), "want response %q, got %q")
	assertSameString(t, "Please type `sector 1`", stp.FailMessage(), "want response %q, got %q")
}

func TestStoryCurrentStep(t *testing.T) {
	str := mesproc.NewStory().Add(mesproc.NewStep().Expect("wow").Respond("yes!"))

	assertSameString(t, "wow", str.Step().Expectation(), "want current step expectation %q, got %q")
	assertSameString(t, "yes!", str.Step().Response(), "want current step expectation %q, got %q")
}

func TestSteps(t *testing.T) {
	str := mesproc.NewStory()

	stp1 := mesproc.NewStep().
		Expect("sector 1").
		Respond("The story of this sector").
		Fail("Please type `sector 1`")
	stp2 := mesproc.NewStep().
		Expect("lulz").Respond("Chilling lulz").Fail("I want to hear `lulz`")

	str.Add(stp1).Add(stp2)

	// TODO: RespondTo has side effects - winding the current step. Should rethink design?
	assertSameString(t, stp1.FailMessage(), str.RespondTo("smth else"), "want response %q, got %q")
	assertSameString(t, stp1.Response(), str.RespondTo("sector 1"), "want response %q, got %q")

	assertSameString(t, stp2.FailMessage(), str.RespondTo("sector 1"), "want response %q, got %q")
	assertSameString(t, stp2.Response(), str.RespondTo("lulz"), "want response %q, got %q")
}

func TestStoryMustLoop(t *testing.T) {
	str := mesproc.NewStory().
		Add(mesproc.NewStep().Expect("s1").Respond("step 1")).
		Add(mesproc.NewStep().Expect("s2").Respond("step 2")).
		Add(mesproc.NewStep().Expect("s3").Respond("step 3"))

	str.RespondTo("s1")
	str.RespondTo("s2")
	str.RespondTo("s3")

	assertSameString(t, "step 1", str.RespondTo("s1"), "want response %q got %q")
}

func TestStoryI18N(t *testing.T) {
	str := mesproc.NewStory().
		I18n(mesproc.I18nMap{
			"ru": {
				mesproc.I18nLanguageChanged: "Язык изменен на русский",
			},
			"kk": {
				mesproc.I18nLanguageChanged: "Язык изменен на казахский",
			},
		})

	assertSameString(t, "Язык изменен на русский", str.RespondTo("/ru"), "want language message %q, got %q")
	assertSameString(t, "Язык изменен на казахский", str.RespondTo("/kk"), "want language message %q, got %q")
	assertSameString(t, mesproc.I18nLanguageChanged, str.RespondTo("/en"), "want language message %q, got %q")
}

func TestRepeatPrevStepWhenChangingLanguage(t *testing.T) {

}
