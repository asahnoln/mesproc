package mesproc

type Story struct {
	steps   []*Step
	curStep int
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
	if s.curStep == len(s.steps) {
		s.curStep = 0
	}

	return s.steps[s.curStep]
}

func (s *Story) RespondTo(m string) string {
	step := s.steps[s.curStep]
	if step.expectation == m {
		s.curStep++
		return step.response
	}

	return step.failMessage
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
