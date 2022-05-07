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

func TestLoadingFromJSON(t *testing.T) {
	// TODO: Should remove folder if exists or better use temp dir
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

	t.Run("Unordered steps", func(t *testing.T) {
		assert.Equal(t, "out of order", str.ResponsesWithLangStepTo(66, "", "unordered")[0].Text(), "want response message to command")
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
