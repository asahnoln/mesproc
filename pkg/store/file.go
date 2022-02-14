package store

import (
	"os"
	"time"
)

type File struct {
	path string
}

func NewFile(path string) *File {
	return &File{path}
}

func (f *File) Save(m string) error {
	file, err := os.CreateTemp(f.path, "review-"+time.Now().String())
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(m))
	if err != nil {
		return err
	}

	return file.Close()
}
