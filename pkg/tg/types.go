package tg

// Update is an object sent by Bot when it receives a message from user
type Update struct {
	Message Message
}

// Chat is a subobject with chat information
type Chat struct {
	ID int
}

type From struct {
	LanguageCode string `json:"language_code"`
}

// Message is a subobject of Update object with info on received message
type Message struct {
	Chat     Chat
	Text     string
	Location *Location
	From     From
}

// Location is a subobject of Message object with info on sent geolocation
type Location struct {
	Longitude, Latitude float64
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

// SendPhoto is an object used to send an audio to a bot
type SendPhoto struct {
	ChatID int    `json:"chat_id"`
	Photo  string `json:"photo"`
}

// SendChatAction is an object used to send a chat action to a bot
type SendChatAction struct {
	ChatID int    `json:"chat_id"`
	Action string `json:"action"`
}

// SetChatID sets chat ID for current sender
func (s *SendPhoto) SetChatID(i int) {
	s.ChatID = i
}

// SetContent sets content for current sender
func (s *SendPhoto) SetContent(a string) {
	s.Photo = a
}

// URL returns Telegram endpoint to process current sender
func (s *SendPhoto) URL() string {
	return "/sendPhoto"
}

func (s *SendPhoto) GetChatID() int {
	return s.ChatID
}

func (s *SendPhoto) ChatAction() string {
	return "upload_photo"
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

func (s *SendAudio) GetChatID() int {
	return s.ChatID
}

func (s *SendAudio) ChatAction() string {
	return "upload_document"
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
