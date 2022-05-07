package story

import (
	"strings"
)

// Response is a struct which story returns in response to a message
type Response struct {
	Additional map[string]interface{}

	original, text, lang string
	shouldAdvance        bool
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
	unordered    map[string]*Step
	curStepIndex int
	i18n         I18nMap
	lang         string
}

// New creates a new Story
func New() *Story {
	return &Story{
		cmds:      make(map[string]*Step),
		unordered: make(map[string]*Step),
	}
}

// Add adds a Step to the Story
func (s *Story) Add(step *Step) *Story {
	s.steps = append(s.steps, step)
	return s
}

// AddCommand adds a command to the story
func (s *Story) AddCommand(step *Step) *Story {
	s.cmds[step.Expectation()] = step
	return s
}

func (s *Story) AddUnordered(step *Step) *Story {
	s.unordered[step.Expectation()] = step
	return s
}

// ResponsesWithLangStepTo return multiple responses from a step with ones
func (s *Story) ResponsesWithLangStepTo(stp int, lang string, m string) []Response {
	rs, l, ok := s.parseAndRespond(stp, lang, m)
	result := make([]Response, len(rs))
	for i, r := range rs {
		result[i] = Response{
			original:      r,
			text:          s.i18n.Line(r, l),
			shouldAdvance: ok,
			lang:          l,
		}

		if len(s.steps) > 0 {
			result[i].Additional = s.steps[s.rotateStep(stp)].additional[i]
		}
	}
	return result
}

func (s *Story) parseAndRespond(stp int, lang string, m string) ([]string, string, bool) {
	if lang == "" {
		lang = "en"
	}

	if r, l, ok := s.parseUnordered(m); ok {
		if l != "" {
			lang = l
		}
		return r, lang, false
	}

	r, ok := s.stepResponsesOrFail(m, lang, s.rotateStep(stp))
	return r, lang, ok
}

// I18n sets i18n localzation for the story
func (s *Story) I18n(i I18nMap) *Story {
	s.i18n = i
	return s
}

func (s *Story) I18nMap() I18nMap {
	return s.i18n
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
	if stp.store != nil {
		err := stp.store.Save(m)
		return err == nil
	}

	if !stp.isGeo {
		return strings.EqualFold(s.i18n.Line(stp.expectation, lang), m)
	}

	return stp.checkGeo(m)
}

func (s *Story) parseUnordered(m string) ([]string, string, bool) {
	lookUp := s.unordered

	if strings.HasPrefix(m, "/") {
		lookUp = s.cmds
		m = m[1:]
	}

	if r, lang, ok := s.processI18nCommand(m); ok {
		return []string{r}, lang, true
	}

	if stp, ok := lookUp[m]; ok {
		return stp.Responses(), "", true
	}

	return nil, "", false
}

func (s *Story) processI18nCommand(c string) (string, string, bool) {
	if _, ok := s.i18n[c]; c == "en" || ok {
		return I18nLanguageChanged, c, true
	}

	return "", "", false
}
