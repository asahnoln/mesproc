package tg

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
	Chat     Chat
	Text     string
	Location *Location
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
