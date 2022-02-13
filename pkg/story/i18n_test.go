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
}

func TestI18nLoadingError(t *testing.T) {
	_, err := story.LoadI18n(strings.NewReader(""))
	require.Error(t, err, "want error when nothing to load from")
}
