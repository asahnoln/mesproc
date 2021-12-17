package mesproc

type Handler interface {
	Handle()
}

func HandleRequests(h Handler) {
	h.Handle()
}
