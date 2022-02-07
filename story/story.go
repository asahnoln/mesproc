package story

import (
	"strings"
)

const (
	// I18nLanguageChanged is a default message returned by ResponseTo if language is changed
	I18nLanguageChanged = "Language changed"
)

// I18nMap holds information on internationalization for the Story.
// It uses English by default as indexes to find appropriate translations in other languages.
type I18nMap map[string]map[string]string

// Story holds information on the current story.
// It has steps and i18n. The story is unfolded by using RespondTo, which always
// advances internal counter to the next step until all steps are processed.
// Then, it starts from the beginning.
type Story struct {
	steps        []*Step
	cmds         map[string]*Step
	curStepIndex int
	i18n         I18nMap
	lang         string
}

// New creates a new Story
func New() *Story {
	return &Story{
		cmds: make(map[string]*Step),
	}
}

// Add adds a Step to the Story
func (s *Story) Add(step *Step) *Story {
	s.steps = append(s.steps, step)
	return s
}

func (s *Story) AddCommand(step *Step) *Story {
	s.cmds[step.Expectation()] = step
	return s
}

// Step return the current step of the story.
func (s *Story) Step() *Step {
	return s.steps[s.curStepIndex]
}

// RespondTo returns a response to given message from the current step.
// Internally, on success, the Story advances internal counter to the next step,
// changing current step.
func (s *Story) RespondTo(m string) string {
	r, _, ok := s.parseAndRespond(s.curStepIndex, s.Language(), m)
	if ok {
		s.curStepIndex = s.rotateStep(s.curStepIndex + 1)
	}

	return r
}

type Response struct {
	text, lang    string
	shouldAdvance bool
}

func (r Response) Text() string {
	return r.text
}

func (r Response) ShouldAdvance() bool {
	return r.shouldAdvance
}

func (r Response) Lang() string {
	return r.lang
}

// RespondWithStepTo returns response from a step indicated by given index
func (s *Story) RespondWithStepTo(stp int, m string) Response {
	r, _, ok := s.parseAndRespond(stp, s.Language(), m)
	return Response{
		text:          r,
		shouldAdvance: ok,
	}
}

func (s *Story) RespondWithLangStepTo(stp int, lang string, m string) Response {
	r, l, ok := s.parseAndRespond(stp, lang, m)
	return Response{
		text:          r,
		shouldAdvance: ok,
		lang:          l,
	}
}

func (s *Story) parseAndRespond(stp int, lang string, m string) (string, string, bool) {
	if r, l, ok := s.parseCommand(m); ok {
		if l != "" {
			lang = l
		}
		return s.getI18nLine(lang, r), lang, false
	}

	r, ok := s.stepResponseOrFail(m, lang, s.rotateStep(stp))
	return s.getI18nLine(lang, r), lang, ok
}

// I18n sets i18n localzation for the story
func (s *Story) I18n(i I18nMap) *Story {
	s.i18n = i
	return s
}

// Language retutns currently set language
func (s *Story) Language() string {
	return s.lang
}

// SetLanguage sets current language of the story according to i18n it has
func (s *Story) SetLanguage(l string) *Story {
	s.lang = l
	return s
}

func (s *Story) rotateStep(stp int) int {
	return stp % len(s.steps)
}

func (s *Story) stepResponseOrFail(m, lang string, stp int) (string, bool) {
	step := s.steps[stp]

	if s.isExpectationCorrect(m, lang, step) {
		return step.response, true
	}

	return step.failMessage, false
}

func (s *Story) isExpectationCorrect(m, lang string, stp *Step) bool {
	if !stp.isGeo {
		return s.getI18nLine(lang, stp.expectation) == m
	}

	return stp.checkGeo(m)
}

func (s *Story) getI18nLine(lang, l string) string {
	if lang != "" {
		if r, ok := s.i18n[lang][l]; ok {
			l = r
		}
	}

	return l
}

func (s *Story) parseCommand(m string) (string, string, bool) {
	if !strings.HasPrefix(m, "/") {
		return "", "", false
	}

	c := m[1:]
	if r, lang, ok := s.processI18nCommand(c); ok {
		return r, lang, true
	}

	if stp, ok := s.cmds[c]; ok {
		return stp.Response(), "", true
	}

	return "", "", false
}

func (s *Story) processI18nCommand(c string) (string, string, bool) {
	if _, ok := s.i18n[c]; c == "en" || ok {
		s.SetLanguage(c)
		return I18nLanguageChanged, c, true
	}

	return "", "", false
}
