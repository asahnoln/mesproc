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
	usrCfgs map[int]usrCfg
	lgr     *log.Logger
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
		usrCfgs: make(map[int]usrCfg),
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

// send sends back a Sender
func (h *Handler) send(u Update) {
	id := u.Message.Chat.ID
	usrCfg := h.usrCfgs[id]
	if usrCfg.lang == "" {
		usrCfg.lang = u.Message.From.LanguageCode
	}

	rs := h.str.ResponsesWithLangStepTo(usrCfg.step, usrCfg.lang, convertText(u))
	rs, translated := h.translateLastResponses(usrCfg, rs)

	for _, r := range rs {
		if t, ok := r.Additional["time"]; ok {
			go func() {
				time.Sleep(t.(time.Duration))
				h.sendResponse(r, id)
			}()
		} else {
			h.sendResponse(r, id)
		}

	}

	usrCfg.lastRs = rs
	h.updateUsrCfg(id, usrCfg, rs[0], translated)
}

func (h *Handler) sendResponse(r story.Response, id int) {
	v := figureSenderType(r.Text())
	v.SetChatID(id)

	h.before(v)

	// TODO: Handle error
	m, _ := json.Marshal(v)
	http.Post(h.target+v.URL(), "application/json", bytes.NewReader(m))
}

func (h *Handler) translateLastResponses(u usrCfg, rs []story.Response) ([]story.Response, bool) {
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
}

func (h *Handler) updateUsrCfg(id int, u usrCfg, r story.Response, translated bool) {
	if r.ShouldAdvance() && !translated {
		u.step++
	}
	u.lang = r.Lang()
	h.usrCfgs[id] = u
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
