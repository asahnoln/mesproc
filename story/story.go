package story

import "strings"

const (
	I18nLanguageChanged = "Language changed"
)

type I18nMap map[string]map[string]string

type Story struct {
	steps   []*Step
	curStep int
	i18n    I18nMap
	lang    string
}

func New() *Story {
	return &Story{}
}

func (s *Story) Add(step *Step) *Story {
	s.steps = append(s.steps, step)
	return s
}

func (s *Story) Step() *Step {
	return s.steps[s.curStep]
}

func (s *Story) RespondTo(m string) string {
	if response, ok := s.parseI18nCommand(m); ok {
		return response
	}

	r, ok := s.stepResponseOrFail(m, s.curStep)
	if ok {
		s.curStep = s.rotateStep(s.curStep + 1)
	}

	return r
}

func (s *Story) RespondWithStepTo(stp int, m string) string {
	r, _ := s.stepResponseOrFail(m, s.rotateStep(stp))
	return r
}

func (s *Story) I18n(i I18nMap) *Story {
	s.i18n = i
	return s
}

func (s *Story) Language() string {
	return s.lang
}

func (s *Story) SetLanguage(l string) *Story {
	s.lang = l
	return s
}

func (s *Story) rotateCurrentStep() {
}

func (s *Story) rotateStep(stp int) int {
	return stp % len(s.steps)
}

func (s *Story) stepResponseOrFail(m string, stp int) (string, bool) {
	var response string
	var ok bool
	step := s.steps[stp]

	if s.isExpectationCorrect(m, step) {
		ok = true
		response = step.response
	} else {
		response = step.failMessage
	}

	return s.getI18nLine(response), ok
}

func (s *Story) isExpectationCorrect(m string, stp *Step) bool {
	if !s.Step().isGeo {
		return s.getI18nLine(stp.expectation) == m
	}

	return s.Step().checkGeo(m)
}

func (s *Story) getI18nLine(l string) string {
	if s.lang != "" {
		r, ok := s.i18n[s.lang][l]
		if ok {
			l = r
		}
	}

	return l
}

func (s *Story) parseI18nCommand(m string) (string, bool) {
	if !strings.HasPrefix(m, "/") {
		return "", false
	}

	c := m[1:]
	if c == "en" {
		return I18nLanguageChanged, true
	}

	if lines, ok := s.i18n[c]; ok {
		s.SetLanguage(c)
		return lines[I18nLanguageChanged], true
	}

	return "", false
}
