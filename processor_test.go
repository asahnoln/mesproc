package mesproc_test

import (
	"testing"

	"github.com/asahnoln/mesproc"
)

func TestReceive(t *testing.T) {
	service := &stubService{}
	mesproc.HandleRequests(service)

	if !service.handled {
		t.Error("service was not handled by HandleRequests")
	}
}

type stubService struct {
	handled bool
}

func (s *stubService) Handle() {
	s.handled = true
}

func TestHttpHandlerHandled(t *testing.T) {
	h := mesproc.NewHttpHandler()
	mesproc.HandleRequests(h)
}
