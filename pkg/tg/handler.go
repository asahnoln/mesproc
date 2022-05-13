// Package tg implements Telegram API for handling given stories.
package tg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/asahnoln/mesproc/pkg/story"
)

const (
	// PrefixAudio identifies text as a sendAudio candidate
	PrefixAudio = "audio:"
	// PrefixPhoto identifies text as a sendPhoto candidate
	PrefixPhoto = "photo:"
)

type usrCfg struct {
	step   int
	lang   string
	lastRs []story.Response
}

// Handler is a Telegram handler, which implements receiving messages from a bot and sending them back
type Handler struct {
	target  string
	str     *story.Story
	usrCfgs map[int]*usrCfg
	lgr     *log.Logger
	timers  []*time.Timer
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
func New(target string, str *story.Story, logger *log.Logger) *Handler {
	return &Handler{
		target:  target,
		str:     str,
		usrCfgs: make(map[int]*usrCfg),
		lgr:     logger,
	}
}

// receive gets an Update from a bot
func (h *Handler) receive(w http.ResponseWriter, r *http.Request) Update {
	var u Update
	// TODO: Handle error
	json.NewDecoder(r.Body).Decode(&u)

	h.logIncoming(u)
	return u
}

func (h *Handler) logIncoming(u Update) {
	if h.lgr != nil {
		h.lgr.Printf("%s: telegram update: %#v", time.Now().Format(time.RFC3339), u)
	}
}

func (h *Handler) logSending(r *http.Response, err error) {
	if h.lgr != nil {
		h.lgr.Printf("%s: response from telegram: %#v, error %#v", time.Now().Format(time.RFC3339), r, err)
	}
}

// send sends back a Sender
func (h *Handler) send(u Update) {
	id := u.Message.Chat.ID
	uCfg := h.prepareUserConfig(id, u)

	if h.runTimedResponses() {
		return
	}

	rs := h.str.ResponsesWithLangStepTo(uCfg.step, uCfg.lang, convertText(u))
	rs, translated := h.translateLastResponses(uCfg, rs)

	for _, r := range rs {
		if t, ok := r.Additional["time"]; ok {
			h.addTimedResponse(r, t.(time.Duration), id)
		} else {
			h.sendResponse(r, id)
		}

	}

	uCfg.lastRs = rs
	h.updateUsrCfg(id, uCfg, rs[0], translated)
}

func (h *Handler) prepareUserConfig(id int, u Update) *usrCfg {
	uCfg, ok := h.usrCfgs[id]
	if !ok {
		h.usrCfgs[id] = &usrCfg{}
		uCfg = h.usrCfgs[id]
	}
	if uCfg.lang == "" {
		uCfg.lang = u.Message.From.LanguageCode
	}
	if u.Message.Text == "/start" {
		uCfg.step = 0
	}

	return uCfg
}

func (h *Handler) addTimedResponse(r story.Response, t time.Duration, id int) {
	timer := time.AfterFunc(t, func() {
		h.sendResponse(r, id)
		// TODO: Write test for this one
		if len(h.timers) > 0 {
			h.timers = h.timers[1:]
		}
	})
	h.timers = append(h.timers, timer)
}

func (h *Handler) runTimedResponses() bool {
	if len(h.timers) > 0 {
		timer := h.timers[0]
		timer.Reset(0)
		return true
	}

	return false
}

func (h *Handler) sendResponse(r story.Response, id int) {
	v := figureSenderType(r.Text())
	v.SetChatID(id)

	h.before(v)

	// TODO: Handle error
	m, _ := json.Marshal(v)
	i := 0
	for {
		resp, err := http.Post(h.target+v.URL(), "application/json", bytes.NewReader(m))
		if resp.StatusCode == http.StatusOK && err == nil {
			break
		}
		if i < 1 {
			h.logSending(resp, err)
		}
		i++
	}
}

func (h *Handler) translateLastResponses(u *usrCfg, rs []story.Response) ([]story.Response, bool) {
	if u.lastRs != nil && rs[0].Lang() != u.lang {
		return h.str.I18nMap().Translate(u.lastRs, rs[0].Lang()), true
	}
	return rs, false
}

func (h *Handler) before(v Sender) {
	if a, ok := v.(*SendAudio); ok {
		m, _ := json.Marshal(SendChatAction{
			ChatID: a.ChatID,
			Action: "upload_document",
		})
		http.Post(h.target+"/sendChatAction", "application/json", bytes.NewReader(m))
	}
	if a, ok := v.(*SendPhoto); ok {
		m, _ := json.Marshal(SendChatAction{
			ChatID: a.ChatID,
			Action: "upload_photo",
		})
		http.Post(h.target+"/sendChatAction", "application/json", bytes.NewReader(m))
	}
}

func (h *Handler) updateUsrCfg(id int, u *usrCfg, r story.Response, translated bool) {
	if r.ShouldAdvance() && !translated {
		u.step++
	}
	u.lang = r.Lang()
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
	switch {
	case strings.HasPrefix(text, PrefixAudio):
		v = &SendAudio{}
		text = text[len(PrefixAudio):]
	case strings.HasPrefix(text, PrefixPhoto):
		v = &SendPhoto{}
		text = text[len(PrefixPhoto):]
	}

	v.SetContent(text)

	return v
}
