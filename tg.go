package mesproc

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
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

type TgSendAudio struct {
	ChatID int    `json:"chat_id"`
	Audio  string `json:"audio"`
}

type TgHandler struct {
	target     string
	lastChatID int
	str        *Story
}

type TgSender interface {
	SetChatID(int)
	SetContent(string)
	URL() string
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
	v := figureSenderType(h.str.RespondTo(message))
	v.SetChatID(h.lastChatID)

	m, _ := json.Marshal(v)

	http.Post(h.target+v.URL(), "application/json", bytes.NewReader(m))
}

func figureSenderType(text string) TgSender {
	var v TgSender = &TgSendMessage{}
	if strings.HasPrefix(text, "audio:") {
		v = &TgSendAudio{}
		text = text[6:]
	}

	v.SetContent(text)

	return v
}

func (s *TgSendAudio) SetChatID(i int) {
	s.ChatID = i
}

func (s *TgSendAudio) SetContent(a string) {
	s.Audio = a
}

func (s *TgSendAudio) URL() string {
	return "/sendAudio"
}

func (s *TgSendMessage) SetChatID(i int) {
	s.ChatID = i
}

func (s *TgSendMessage) SetContent(a string) {
	s.Text = a
}

func (s *TgSendMessage) URL() string {
	return "/sendMessage"
}
