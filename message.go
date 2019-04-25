package slackbot

// Message represents a message.
type Message struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	SubType string `json:"subtype"`
	Channel string `json:"channel"`
	UserID  string `json:"user"`
	Text    string `json:"text"`
	Time    int64  `json:"time"`
}
