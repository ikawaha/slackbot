package webapi

// MessageResponse represents the response of the chat.postMessage API.
type MessageResponse struct {
	OK      bool    `json:"ok,omitempty"`
	Error   string  `json:"error,omitempty"`
	Channel string  `json:"channel,omitempty"`
	TS      string  `json:"ts,omitempty"`
	Message Message `json:"message,omitempty"`
}

// Message represents the Slack message.
type Message struct {
	Text        string       `json:"text,omitempty"`
	Username    string       `json:"username,omitempty"`
	BotID       string       `json:"bot_id,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	Type        string       `json:"type,omitempty"`
	SubType     string       `json:"sub_type,omitempty"`
	TS          string       `json:"ts,omitempty"`
}

// Attachment is a part of the Message.
type Attachment struct {
	Text     string `json:"text,omitempty"`
	ID       int    `json:"id,omitempty"`
	Fallback string `json:"fallback,omitempty"`
}
