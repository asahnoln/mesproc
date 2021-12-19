package mesproc

type AnswerMap map[string]string

type Story struct {
	m AnswerMap
}

func NewStory(m AnswerMap) *Story {
	return &Story{m}
}

func (s *Story) Respond(m string) string {
	return s.m[m]
}
