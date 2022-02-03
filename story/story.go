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

	s.checkCurrentStep()
	return s.stepResponseOrFail(m, s.curStep, true)
}

func (s *Story) RespondWithStepTo(stp int, m string) string {
	return s.stepResponseOrFail(m, stp, false)
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

func (s *Story) checkCurrentStep() {
	if s.curStep == len(s.steps) {
		s.curStep = 0
	}
}

func (s *Story) stepResponseOrFail(m string, stp int, incCurStep bool) string {
	var response string
	step := s.steps[stp]
	if s.isExpectationCorrect(m, step) {
		response = step.response
		if incCurStep {
			s.curStep++
		}
	} else {
		response = step.failMessage
	}

	return s.getI18nLine(response)
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
