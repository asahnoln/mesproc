// Package tg implements Telegram API for handling given stories.
package tg

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/asahnoln/mesproc/story"
)

const (
	PREFIX_AUDIO = "audio:" // PREFIX_AUDIO identifies text as a sendAudio candidate
)

// Update is an object sent by Bot when it receives a message from user
type Update struct {
	Message Message
}

// Chat is a subobject with chat information
type Chat struct {
	ID int
}

// Message is a subobject of Update object with info on received message
type Message struct {
	Chat Chat
	Text string
}

// SendMessage is an object used to send a message to a bot
type SendMessage struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

// SendAudio is an object used to send an audio to a bot
type SendAudio struct {
	ChatID int    `json:"chat_id"`
	Audio  string `json:"audio"`
}

// Handler is a Telegram handler, which implements receiving messages from a bot and sending them back
type Handler struct {
	target string
	str    *story.Story
}

// Sender is an interface for different sending options, like sendMessage, sendAudio etc.
// SetChatID and SetContent are used to set internal fields with correct values to send them
// back to Telegram Bot API.
// URL is used to receive correct endpoint for parcticular Sender.
type Sender interface {
	SetChatID(int)
	SetContent(string)
	URL() string
}

// New creates a Telegram handler.
func New(target string, str *story.Story) *Handler {
	return &Handler{target: target, str: str}
}

// receive gets an Update from a bot
func (h *Handler) receive(w http.ResponseWriter, r *http.Request) Update {
	var u Update
	_ = json.NewDecoder(r.Body).Decode(&u)
	return u
}

// send sends back a Sender
func (h *Handler) send(u Update) {
	v := figureSenderType(h.str.RespondTo(u.Message.Text))
	v.SetChatID(u.Message.Chat.ID)

	m, _ := json.Marshal(v)

	http.Post(h.target+v.URL(), "application/json", bytes.NewReader(m))
}

// ServeHTTP implements http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.send(h.receive(w, r))
}

// figureSenderType uses received text as a way to figure out what should be sent back
func figureSenderType(text string) Sender {
	var v Sender = &SendMessage{}
	if strings.HasPrefix(text, PREFIX_AUDIO) {
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
