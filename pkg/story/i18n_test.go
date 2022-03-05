package story_test

import (
	"os"
	"strings"
	"testing"

	"github.com/asahnoln/mesproc/pkg/story"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestI18nLoadingFromJSON(t *testing.T) {
	f, err := os.Open("testdata/i18n.json")
	require.NoError(t, err, "unexepcted error while loading test i18n json")
	defer f.Close()

	i18n, err := story.LoadI18n(f)

	require.NoError(t, err, "unexpected error while loading file into i18n")

	assert.Equal(t, "шаг 1", i18n.Line("step 1", "ru"), "want russian translation")
	assert.Equal(t, "doesn't exist", i18n.Line("doesn't exist", "ru"), "want default translation")
}

func TestI18nLoadingError(t *testing.T) {
	_, err := story.LoadI18n(strings.NewReader(""))
	require.Error(t, err, "want error when nothing to load from")
}

func TestI18nTranslateResponses(t *testing.T) {
	str := story.New().
		AddCommand(story.NewStep().Expect("command").Respond("cmdResponse")).
		Add(story.NewStep().Expect("message").Respond("one", "two").Fail("fail")).
		I18n(story.I18nMap{
			"ru": {
				"cmdResponse": "кмдОтвет",
				"one":         "один",
				"two":         "два",
			},
		})

	t.Run("common reply", func(t *testing.T) {
		rs := str.ResponsesWithLangStepTo(0, "ru", "message")

		localized := str.I18nMap().Translate(rs, "en")

		require.Len(t, localized, 2, "want same count of translated responses")
		assert.Equal(t, "en", localized[0].Lang(), "want Russian language responses")
		assert.Equal(t, "one", localized[0].Text())
	})

	t.Run("command", func(t *testing.T) {
		rs := str.ResponsesWithLangStepTo(99, "ru", "/command")
		localized := str.I18nMap().Translate(rs, "en")
		assert.Equal(t, "cmdResponse", localized[0].Text())
	})
}
