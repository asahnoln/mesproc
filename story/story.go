package story

import (
	"strings"
)

// Response is a struct which story returns in response to a message
type Response struct {
	text, lang    string
	shouldAdvance bool
}

// Text returns text of response
func (r Response) Text() string {
	return r.text
}

// ShouldAdvance returns whether the step index should advance in the story for the next response
func (r Response) ShouldAdvance() bool {
	return r.shouldAdvance
}

// Lang returns the language of the response
func (r Response) Lang() string {
	return r.lang
}

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

	return r[0]
}

// RespondWithStepTo returns response from a step indicated by given index
func (s *Story) RespondWithStepTo(stp int, m string) Response {
	r, _, ok := s.parseAndRespond(stp, s.Language(), m)
	return Response{
		text:          r[0],
		shouldAdvance: ok,
	}
}

func (s *Story) RespondWithLangStepTo(stp int, lang string, m string) Response {
	r, l, ok := s.parseAndRespond(stp, lang, m)
	return Response{
		text:          r[0],
		shouldAdvance: ok,
		lang:          l,
	}
}

func (s *Story) ResponsesWithLangStepTo(stp int, lang string, m string) []Response {
	rs, l, ok := s.parseAndRespond(stp, lang, m)
	result := make([]Response, len(rs))
	for i, r := range rs {
		result[i] = Response{
			text:          r,
			shouldAdvance: ok,
			lang:          l,
		}
	}
	return result
}

func (s *Story) parseAndRespond(stp int, lang string, m string) ([]string, string, bool) {
	if r, l, ok := s.parseCommand(m); ok {
		if l != "" {
			lang = l
		}
		return s.getI18nLines(lang, r), lang, false
	}

	r, ok := s.stepResponsesOrFail(m, lang, s.rotateStep(stp))
	return s.getI18nLines(lang, r), lang, ok
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

func (s *Story) stepResponsesOrFail(m, lang string, stp int) ([]string, bool) {
	step := s.steps[stp]

	if s.isExpectationCorrect(m, lang, step) {
		return step.Responses(), true
	}

	return []string{step.failMessage}, false
}

func (s *Story) isExpectationCorrect(m, lang string, stp *Step) bool {
	if !stp.isGeo {
		return strings.EqualFold(s.getI18nLine(lang, stp.expectation), m)
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

func (s *Story) getI18nLines(lang string, ls []string) []string {
	for i, l := range ls {
		ls[i] = s.getI18nLine(lang, l)
	}

	return ls
}

func (s *Story) parseCommand(m string) ([]string, string, bool) {
	if !strings.HasPrefix(m, "/") {
		return nil, "", false
	}

	c := m[1:]
	if r, lang, ok := s.processI18nCommand(c); ok {
		return []string{r}, lang, true
	}

	if stp, ok := s.cmds[c]; ok {
		return stp.Responses(), "", true
	}

	return nil, "", false
}

func (s *Story) processI18nCommand(c string) (string, string, bool) {
	if _, ok := s.i18n[c]; c == "en" || ok {
		s.SetLanguage(c)
		return I18nLanguageChanged, c, true
	}

	return "", "", false
}
