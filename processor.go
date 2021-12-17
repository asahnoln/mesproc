package mesproc

import (
	"net/http"
)

// Server is an http server which uses a Handler to Receive a request, handle it and Send a response
type Server struct {
	h Handler
}

// Handler is an object which can Receive an http reqiest, handle it with Request method,
// and then use its contents to Send a response back to the service
type Handler interface {
	Receive(http.ResponseWriter, *http.Request) string
	Send(string)
}

// NewServer creates a Server with given Handler
func NewServer(h Handler) *Server {
	return &Server{
		h,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.h.Send(s.h.Receive(w, r))
}
