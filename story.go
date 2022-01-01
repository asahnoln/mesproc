package mesproc

import "strings"

const (
	I18nLanguageChanged = "Language changed to English"
)

type I18nMap map[string]map[string]string

type Story struct {
	steps   []*Step
	curStep int
	i18n    I18nMap
}

type Step struct {
	expectation string
	response    string
	failMessage string
}

func NewStory() *Story {
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
	return s.stepResponseOrFail(m)
}

func (s *Story) I18n(i I18nMap) *Story {
	s.i18n = i
	return s
}

func (s *Story) checkCurrentStep() {
	if s.curStep == len(s.steps) {
		s.curStep = 0
	}
}

func (s *Story) stepResponseOrFail(m string) string {
	step := s.steps[s.curStep]
	if step.expectation == m {
		s.curStep++
		return step.response
	}

	return step.failMessage
}

func (s *Story) parseI18nCommand(m string) (string, bool) {
	var (
		response string
		ok       bool
	)
	if !strings.HasPrefix(m, "/") {
		return response, ok
	}

	c := m[1:]
	if c == "en" {
		return I18nLanguageChanged, true
	}

	if lines, ok := s.i18n[c]; ok {
		return lines[I18nLanguageChanged], true
	}

	return "", false
}

func NewStep() *Step {
	return &Step{}
}

func (s *Step) Expect(e string) *Step {
	s.expectation = e
	return s
}

func (s *Step) Respond(r string) *Step {
	s.response = r
	return s
}

func (s *Step) Expectation() string {
	return s.expectation
}

func (s *Step) Response() string {
	return s.response
}

func (s *Step) Fail(e string) *Step {
	s.failMessage = e
	return s
}

func (s *Step) FailMessage() string {
	return s.failMessage
}
