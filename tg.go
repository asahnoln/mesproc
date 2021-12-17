package mesproc

import (
	"net/http"
)

type TgSendMessage struct {
	Text string
}

type TgHandler struct {
	target string
}

func NewTgHandler(target string) *TgHandler {
	return &TgHandler{target}
}

func (h *TgHandler) Receive(w http.ResponseWriter, r *http.Request) string {
	return ""
}

func (h *TgHandler) Send(string) {
	_ = TgSendMessage{"LILKI"}
	// json.NewEncoder(w io.Writer)
	http.Post(h.target, "", nil)
}
