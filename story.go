package mesproc

type Story struct {
	steps []*Step
}

type Step struct {
	response string
}

func NewStory() *Story {
	return &Story{}
}

func (s *Story) Add(step *Step) {
	s.steps = append(s.steps, step)
}

func (s *Story) RespondTo(m string) string {
	var r string
	switch m {
	case "sector 1":
		r = s.steps[0].response
	case "lulz":
		r = s.steps[1].response
	}

	return r
}

func NewStep() *Step {
	return &Step{}
}

func (s *Step) Expect(string) *Step {
	return s
}

func (s *Step) Respond(r string) *Step {
	s.response = r
	return s
}

func (s *Step) Response() string {
	return s.response
}
