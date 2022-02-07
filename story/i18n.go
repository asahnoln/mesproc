package story

import (
	"encoding/json"
	"io"
)

const (
	// I18nLanguageChanged is a default message returned by ResponseTo if language is changed
	I18nLanguageChanged = "Language changed"
)

// I18nMap holds information on internationalization for the Story.
// It uses English by default as indexes to find appropriate translations in other languages.
type I18nMap map[string]map[string]string

func LoadI18n(r io.Reader) (I18nMap, error) {
	m := make(I18nMap)
	err := json.NewDecoder(r).Decode(&m)
	return m, err
}

func (m I18nMap) Line(line, lang string) string {
	return m[lang][line]
}
