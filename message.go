package slackbot

// Message represents a slack message.
type Message struct {
	ID        uint64 `json:"id"`
	Type      string `json:"type"`
	SubType   string `json:"subtype"`
	Channel   string `json:"channel"`
	UserID    string `json:"user"`
	Text      string `json:"text"`
	Timestamp string `json:"ts"` // unix timestamp e.g. "1355517523.000005"
	Time      int64  `json:"time"`
}
