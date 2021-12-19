package mesproc

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type TgUpdate struct {
	Message string
}

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
	var u TgUpdate
	_ = json.NewDecoder(r.Body).Decode(&u)
	return u.Message
}

func (h *TgHandler) Send(message string) {
	var text string
	// TODO: Use story module
	switch message {
	case "/ru":
		text = "Выберите сектор"
	case "/en":
		text = "Choose sector"
	}
	m, _ := json.Marshal(TgSendMessage{
		Text: text,
	})
	http.Post(h.target+"/sendMessage", "", bytes.NewReader(m))
}
