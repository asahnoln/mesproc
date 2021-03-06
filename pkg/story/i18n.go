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

// LoadI18n loads i18n from json file to the struct
func LoadI18n(r io.Reader) (I18nMap, error) {
	m := make(I18nMap)
	err := json.NewDecoder(r).Decode(&m)
	return m, err
}

// Line returns a translated line from I18n
func (m I18nMap) Line(line, lang string) string {
	l, ok := m[lang][line]
	if !ok {
		return line
	}

	return l
}

func (m I18nMap) Translate(rs []Response, lang string) []Response {
	result := make([]Response, len(rs))

	for i, r := range rs {
		result[i] = Response{
			original:      r.original,
			text:          m.Line(r.original, lang),
			lang:          lang,
			shouldAdvance: r.shouldAdvance,
		}
	}
	return result
}
