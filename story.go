package mesproc

type AnswerMap map[string]string

type Story struct {
	asnwers    AnswerMap
	gotAnswers map[string]bool
}

func NewStory(m AnswerMap) *Story {
	return &Story{
		asnwers:    m,
		gotAnswers: make(map[string]bool),
	}
}

func (s *Story) Respond(m string) string {
	if a := s.gotAnswers[m]; a {
		switch m {
		case "sector 1":
			return "Please type `lulz`"
		case "lulz":
			return "Please type `winners 5`"
		case "winners 5":
			return "Please type `guds`"
		}
	}
	s.gotAnswers[m] = true
	return s.asnwers[m]
}
