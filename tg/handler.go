package tg

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/asahnoln/mesproc/story"
)

type Update struct {
	Message Message
}

type Chat struct {
	ID int
}

type Message struct {
	Chat Chat
	Text string
}

type SendMessage struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

type SendAudio struct {
	ChatID int    `json:"chat_id"`
	Audio  string `json:"audio"`
}

type Handler struct {
	target     string
	lastChatID int
	str        *story.Story
}

type Sender interface {
	SetChatID(int)
	SetContent(string)
	URL() string
}

func New(target string, str *story.Story) *Handler {
	return &Handler{target: target, str: str}
}

func (h *Handler) receive(w http.ResponseWriter, r *http.Request) Update {
	var u Update
	_ = json.NewDecoder(r.Body).Decode(&u)
	h.lastChatID = u.Message.Chat.ID
	return u
}

func (h *Handler) send(u Update) {
	v := figureSenderType(h.str.RespondTo(u.Message.Text))
	v.SetChatID(h.lastChatID)

	m, _ := json.Marshal(v)

	http.Post(h.target+v.URL(), "application/json", bytes.NewReader(m))
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.send(h.receive(w, r))
}

func figureSenderType(text string) Sender {
	var v Sender = &SendMessage{}
	if strings.HasPrefix(text, "audio:") {
		v = &SendAudio{}
		text = text[6:]
	}

	v.SetContent(text)

	return v
}

func (s *SendAudio) SetChatID(i int) {
	s.ChatID = i
}

func (s *SendAudio) SetContent(a string) {
	s.Audio = a
}

func (s *SendAudio) URL() string {
	return "/sendAudio"
}

func (s *SendMessage) SetChatID(i int) {
	s.ChatID = i
}

func (s *SendMessage) SetContent(a string) {
	s.Text = a
}

func (s *SendMessage) URL() string {
	return "/sendMessage"
}
