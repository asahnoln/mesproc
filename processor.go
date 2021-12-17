package mesproc

import (
	"net/http"
)

type Server struct {
	h Handler
}

type Handler interface {
	Receive(http.ResponseWriter, *http.Request) string
	Send(string)
}

func NewServer(h Handler) *Server {
	return &Server{
		h,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.h.Send(s.h.Receive(w, r))
}
