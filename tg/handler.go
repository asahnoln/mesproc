// Package tg implements Telegram API for handling given stories.
package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/asahnoln/mesproc/story"
)

const (
	// PrefixAudio identifies text as a sendAudio candidate
	PrefixAudio = "audio:"
)

// Handler is a Telegram handler, which implements receiving messages from a bot and sending them back
type Handler struct {
	target   string
	str      *story.Story
	usrSteps map[int]int
}

// Sender is an interface for different sending options, like sendMessage, sendAudio etc.
// SetChatID and SetContent are used to set internal fields with correct values to send them
// back to Telegram Bot API.
// URL is used to receive correct endpoint for parcticular Sender.
type Sender interface {
	SetChatID(int)     // SetChatID sets chat ID for current sender
	SetContent(string) // SetContent sets content for current sender
	URL() string       // URL returns Telegram endpoint to process current sender
}

// New creates a Telegram handler.
func New(target string, str *story.Story) *Handler {
	return &Handler{
		target:   target,
		str:      str,
		usrSteps: make(map[int]int),
	}
}

// receive gets an Update from a bot
func (h *Handler) receive(w http.ResponseWriter, r *http.Request) Update {
	var u Update
	// TODO: Handle error
	json.NewDecoder(r.Body).Decode(&u)
	return u
}

// send sends back a Sender
func (h *Handler) send(u Update) {
	id := u.Message.Chat.ID
	r := h.str.RespondWithStepTo(h.usrSteps[id], convertText(u))
	v := figureSenderType(r.Text())
	v.SetChatID(id)

	// TODO: Handle error
	m, _ := json.Marshal(v)
	http.Post(h.target+v.URL(), "application/json", bytes.NewReader(m))

	if r.ShouldAdvance() {
		h.usrSteps[id]++
	}
}

// ServeHTTP implements http.Handler
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.send(h.receive(w, r))
}

// convertText converts Update info into text usable by Story
func convertText(u Update) string {
	text := u.Message.Text
	if u.Message.Location != nil {
		text = fmt.Sprintf("%f,%f", u.Message.Location.Latitude, u.Message.Location.Longitude)
	}
	return text
}

// figureSenderType uses received text as a way to figure out what should be sent back
func figureSenderType(text string) Sender {
	var v Sender = &SendMessage{}
	if strings.HasPrefix(text, PrefixAudio) {
		v = &SendAudio{}
		text = text[len(PrefixAudio):]
	}

	v.SetContent(text)

	return v
}

// SetChatID sets chat ID for current sender
func (s *SendAudio) SetChatID(i int) {
	s.ChatID = i
}

// SetContent sets content for current sender
func (s *SendAudio) SetContent(a string) {
	s.Audio = a
}

// URL returns Telegram endpoint to process current sender
func (s *SendAudio) URL() string {
	return "/sendAudio"
}

// SetChatID sets chat ID for current sender
func (s *SendMessage) SetChatID(i int) {
	s.ChatID = i
}

// SetContent sets content for current sender
func (s *SendMessage) SetContent(a string) {
	s.Text = a
}

// URL returns Telegram endpoint to process current sender
func (s *SendMessage) URL() string {
	return "/sendMessage"
}
