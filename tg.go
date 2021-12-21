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
	str        *Story
}

func NewTgHandler(target string, str *Story) *TgHandler {
	return &TgHandler{target: target, str: str}
}

func (h *TgHandler) Receive(w http.ResponseWriter, r *http.Request) string {
	var u TgUpdate
	_ = json.NewDecoder(r.Body).Decode(&u)
	h.lastChatID = u.Message.Chat.ID
	return u.Message.Text
}

func (h *TgHandler) Send(message string) {
	m, _ := json.Marshal(TgSendMessage{
		ChatID: h.lastChatID,
		Text:   h.str.RespondTo(message),
	})
	http.Post(h.target+"/sendMessage", "application/json", bytes.NewReader(m))
}
