package story

type Step struct {
	expectation string
	response    string
	failMessage string
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
