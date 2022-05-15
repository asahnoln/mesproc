package store_test

import (
	"os"
	"path"
	"testing"

	"github.com/asahnoln/mesproc/pkg/store"
	"github.com/stretchr/testify/require"
)

func TestFileStore(t *testing.T) {
	dir := createDir(t)
	defer removeDir(t, dir)

	s := store.NewFile(dir)
	require.Implements(t, (*store.Step)(nil), s, "File store must implement Step interface")

	err := s.Save("my review")
	require.NoError(t, err, "unexpected error while saving")

	dirs, err := os.ReadDir(dir)
	require.NoError(t, err, "unexpected error while opening the file")

	data, err := os.ReadFile(path.Join(dir, dirs[0].Name()))
	require.NoError(t, err, "unexpected error while reading the file")
	require.EqualValues(t, "my review", data, "want exact saved content")

}

func TestFileSuccessors(t *testing.T) {
	dir := createDir(t)
	defer removeDir(t, dir)

	s := store.NewFile(dir)
	_ = s.Save("my review 1")
	_ = s.Save("my review 2")
	_ = s.Save("my review 3")

	dirs, err := os.ReadDir(dir)
	require.NoError(t, err, "unexpected error while opening the file")

	data, err := os.ReadFile(path.Join(dir, dirs[2].Name()))
	require.NoError(t, err, "unexpected error while reading the file")
	require.EqualValues(t, "my review 3", data, "want exact saved content")
}

func TestWrongPathFileError(t *testing.T) {
	s := store.NewFile("nowherefound")
	err := s.Save("anything")
	require.Error(t, err, "want saving error")
}

func createDir(t testing.TB) string {
	dir, err := os.MkdirTemp("", "filestore")
	require.NoError(t, err, "unexpected error while creating tmp dir")

	return dir
}

func removeDir(t testing.TB, dir string) {
	require.NoError(t, os.RemoveAll(dir), "unexpected error while removing the test dir")
}
