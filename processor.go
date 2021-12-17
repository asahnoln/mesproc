package mesproc

import (
	"encoding/json"
	"net/http"
	"strings"
)

const (
	OK_ANSWER = "ok"
)

type AnswerMap map[string]string

type Server struct {
	m      AnswerMap
	target string
}

func NewServer(m AnswerMap, target string) *Server {
	return &Server{
		m:      m,
		target: target,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var m struct {
		Message string
	}
	_ = json.NewDecoder(r.Body).Decode(&m)

	http.Post(s.target, "", strings.NewReader(`{"message": "`+s.m[m.Message]+`"}`))

	w.Write([]byte(OK_ANSWER))
}
