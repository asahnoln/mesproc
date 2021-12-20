package mesproc

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type TgUpdate struct {
	Message TgMessage
}

type TgChat struct {
	ID int
}

type TgMessage struct {
	Chat TgChat
	Text string
}

type TgSendMessage struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type TgHandler struct {
	target     string
	lastChatID int
}

func NewTgHandler(target string) *TgHandler {
	return &TgHandler{target: target}
}

func (h *TgHandler) Receive(w http.ResponseWriter, r *http.Request) string {
	var u TgUpdate
	_ = json.NewDecoder(r.Body).Decode(&u)
	h.lastChatID = u.Message.Chat.ID
	return u.Message.Text
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
		ChatID: h.lastChatID,
		Text:   text,
	})
	http.Post(h.target+"/sendMessage", "application/json", bytes.NewReader(m))
}
