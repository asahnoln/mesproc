package story

import "strings"

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
	curStepIndex int
	i18n         I18nMap
	lang         string
}

// New creates a new Story
func New() *Story {
	return &Story{}
}

// Add adds a Step to the Story
func (s *Story) Add(step *Step) *Story {
	s.steps = append(s.steps, step)
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
	if response, ok := s.parseI18nCommand(m); ok {
		return response
	}

	r, ok := s.stepResponseOrFail(m, s.curStepIndex)
	if ok {
		s.curStepIndex = s.rotateStep(s.curStepIndex + 1)
	}

	return r
}

// RespondWithStepTo returns response from a step indicated by given index
func (s *Story) RespondWithStepTo(stp int, m string) string {
	r, _ := s.stepResponseOrFail(m, s.rotateStep(stp))
	return r
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
